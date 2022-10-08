locals {
  services = toset([
    "sqladmin.googleapis.com",
    "run.googleapis.com",
    "secretmanager.googleapis.com",
    "appengine.googleapis.com",
    "iam.googleapis.com",
    "artifactregistry.googleapis.com",
    "iamcredentials.googleapis.com"
  ])
}

resource "google_project_service" "service" {
  for_each = local.services
  project  = var.project
  service  = each.value
}
