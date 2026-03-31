#!/bin/sh

echo "Waiting for kumo..."
for i in $(seq 1 30); do
  if aws --endpoint-url=http://localhost:4566 s3 ls 2>/dev/null; then
    echo "kumo is ready"
    break
  fi
  echo "Waiting... ($i/30)"
  sleep 2
done

echo "Creating S3 buckets..."

aws --endpoint-url=http://localhost:4566 s3 mb s3://local-private-asset-bucket/
echo 'private bucket created!'

aws --endpoint-url=http://localhost:4566 s3 mb s3://local-public-asset-bucket/
echo 'public bucket created!'

echo "S3 setup Done!"
