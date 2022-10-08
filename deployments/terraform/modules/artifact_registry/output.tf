output "container_registry_path" {
  value = "${var.location}-docker.pkg.dev/${var.project}/${google_artifact_registry_repository.registry.repository_id}"
}
