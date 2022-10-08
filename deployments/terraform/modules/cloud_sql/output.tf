output "db_password" {
  value = random_password.db_password.result
}

output "db_user" {
  value = google_sql_user.app_user
}

output "db_connection_name" {
  value = google_sql_database_instance.instance.connection_name
}

output "db_name" {
  value = google_sql_database.database.name
}
