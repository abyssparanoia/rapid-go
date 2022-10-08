locals {
  backend_roles = [
    "roles/cloudsql.client",
    "roles/storage.admin",
    "roles/firebase.admin",
    "roles/secretmanager.secretAccessor",
    "roles/iam.serviceAccountTokenCreator",
    "roles/pubsub.publisher",
    "roles/pubsub.subscriber",
    "roles/cloudprofiler.agent",
    "roles/cloudkms.signerVerifier"
  ]

  cloud_run_services = {
    api = { name = "api", args = ["http-server", "run"], min_scale = 0, max_scale = 5, },
  }
}