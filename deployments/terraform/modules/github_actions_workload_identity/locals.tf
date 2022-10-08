locals {
  roles = [
    "roles/cloudsql.client",
    "roles/run.admin",
    "roles/iam.serviceAccountUser",
    "roles/storage.admin",
    "roles/artifactregistry.writer"
  ]
}
