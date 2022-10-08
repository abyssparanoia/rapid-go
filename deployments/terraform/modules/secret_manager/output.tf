output "db_password_secret_id" {
  value = google_secret_manager_secret.db_password.secret_id
}

output "fincode_api_key_secret_id" {
  value = google_secret_manager_secret.fincode_api_key.secret_id
}

output "onchain_private_key_secret_id" {
  value = google_secret_manager_secret.onchain_private_key.secret_id
}
