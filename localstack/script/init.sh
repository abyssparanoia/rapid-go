#!/bin/sh

echo "S3 setup start!"
echo "Creating S3 buckets..."

# Create private bucket
aws --endpoint-url=http://aws:4566 s3 mb s3://local-private-asset-bucket/
echo 'private bucket created!'

# Create public bucket
aws --endpoint-url=http://aws:4566 s3 mb s3://local-public-asset-bucket/
echo 'public bucket created!'

echo "S3 setup Done!"