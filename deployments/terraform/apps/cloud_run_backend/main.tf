resource "google_service_account" "backend" {
  account_id   = "app-backend"
  display_name = "Backend Service Account"
}

resource "google_project_iam_member" "backend" {
  for_each = toset(local.backend_roles)
  project  = var.project
  role     = each.value
  member   = "serviceAccount:${google_service_account.backend.email}"
}

resource "google_cloud_run_service" "services" {
  for_each = local.cloud_run_services
  name     = each.value.name
  location = var.location
  project  = var.project

  template {
    spec {
      service_account_name = google_service_account.backend.email

      containers {
        image = "${var.registry_path}/backend:latest"
        ports {
          container_port = 8080
        }
        args = each.value.args

        env {
          name  = "ENV"
          value = var.environment
        }
        env {
          name  = "GCP_PROJECT_ID"
          value = var.project
        }
        env {
          name  = "SERVICE_NAME"
          value = each.value.name
        }
        env {
          name  = "MIN_LOG_SEVERITY"
          value = "DEBUG"
        }
        env {
          name  = "DB_HOST"
          value = "unix(/cloudsql/${var.db_connection_name})"
        }
        env {
          name  = "DB_DATABASE"
          value = var.db_name
        }
        env {
          name  = "DB_USER"
          value = var.db_user
        }
        env {
          name = "DB_PASSWORD"
          value_from {
            secret_key_ref {
              name = var.db_password_secret_id
              key  = var.db_password_secret_version
            }
          }
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/minScale"      = each.value.min_scale
        "autoscaling.knative.dev/maxScale"      = each.value.max_scale
        "run.googleapis.com/cpu-throttling"     = each.value.min_scale == 0 ? "true" : "false"
        "run.googleapis.com/cloudsql-instances" = var.db_connection_name
        "run.googleapis.com/client-name"        = "terraform"
      }
    }
  }

  autogenerate_revision_name = true

  traffic {
    percent         = 100
    latest_revision = true
  }

  lifecycle {
    ignore_changes = [
      template[0].spec[0].containers[0].image,
      template[0].metadata[0].annotations["run.googleapis.com/client-name"],
      template[0].metadata[0].annotations["run.googleapis.com/client-version"],
      template[0].metadata[0].annotations["client.knative.dev/user-image"],
    ]
  }
}

resource "google_cloud_run_service_iam_member" "run_all_users" {
  service  = "api"
  location = var.location
  role     = "roles/run.invoker"
  member   = "allUsers"
  depends_on = [
    google_cloud_run_service.services
  ]
}
