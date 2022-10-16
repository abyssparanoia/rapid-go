output "secret_id" {
  value = google_secret_manager_secret.secret.secret_id
}

output "version" {
  value = google_secret_manager_secret_version.secret.version
}
