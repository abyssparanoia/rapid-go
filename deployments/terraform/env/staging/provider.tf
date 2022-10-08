provider "google" {
  project = var.project
  region  = var.location
}

terraform {
  backend "gcs" {
    bucket = "dev-rapid-go-terraform-state-store"
  }
}
