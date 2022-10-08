resource "google_storage_bucket" "buckets" {
  name          = var.bucket_name
  location      = "asia-northeast1"
  storage_class = "STANDARD"

  cors {
    origin          = ["*"]
    method          = ["GET", "HEAD", "PUT", "POST", "DELETE", "OPTIONS"]
    response_header = ["content-type", "cache-control", "x-requested-with"]
    max_age_seconds = 3600
  }
}
