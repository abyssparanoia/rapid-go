resource "google_artifact_registry_repository" "registry" {
  location      = var.location
  repository_id = var.repository_id
  format        = "DOCKER"
}
