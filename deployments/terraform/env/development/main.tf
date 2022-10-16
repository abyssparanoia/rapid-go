module "gcp_service" {
  source  = "../../modules/gcp_service"
  project = var.project
}

module "github_actions_workload_identity" {
  source = "../../modules/github_actions_workload_identity"

  project      = var.project
  location     = var.location
  repositories = ["abyssparanoia/rapid-go"]

  depends_on = [
    module.gcp_service
  ]
}


module "cloud_sql" {
  source            = "../../modules/cloud_sql"
  location          = var.location
  tier              = local.db_tier
  disk_type         = local.db_disk_type
  availability_type = local.db_availability_type
  db_instance_name  = "master"
  db_name           = "maindb"
  db_user           = "app_user"

  depends_on = [
    module.gcp_service
  ]
}

module "secret_manager_db_password" {
  source    = "../../modules/secret_manager"
  secret_id = "db-passowrd"
  value     = module.cloud_sql.db_password

  depends_on = [
    module.gcp_service
  ]
}

module "artifact_registry" {
  source        = "../../modules/artifact_registry"
  project       = var.project
  location      = var.location
  repository_id = "rapid-go"

  depends_on = [
    module.gcp_service
  ]
}

module "cloud_run_backend" {
  source = "../../apps/cloud_run_backend"

  project                    = var.project
  location                   = var.location
  enviroment                 = local.enviroment
  registry_path              = module.artifact_registry.container_registry_path
  db_connection_name         = module.cloud_sql.db_connection_name
  db_name                    = module.cloud_sql.db_name
  db_user                    = module.cloud_sql.db_user
  db_password_secret_id      = module.secret_manager_db_password.secret_id
  db_password_secret_version = module.secret_manager_db_password.version

  depends_on = [
    module.gcp_service
  ]
}
