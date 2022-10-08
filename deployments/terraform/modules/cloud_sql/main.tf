resource "google_sql_database_instance" "instance" {
  name             = var.db_instance_name
  database_version = "MYSQL_8_0"
  region           = var.location

  settings {
    tier              = var.tier
    disk_type         = var.disk_type
    availability_type = var.availability_type

    database_flags {
      name  = "character_set_server"
      value = "utf8mb4"
    }

    backup_configuration {
      location           = "asia"
      enabled            = true
      binary_log_enabled = true
    }

  }
}

resource "google_sql_database" "database" {
  name     = var.db_name
  instance = google_sql_database_instance.instance.name
}

resource "random_password" "db_password" {
  length = 16
}

resource "google_sql_user" "app_user" {
  name     = var.db_user
  instance = google_sql_database_instance.instance.name
  password = random_password.db_password.result
}
