package gcs

import (
	"context"

	"cloud.google.com/go/storage"
)

func NewClient(ctx context.Context) *storage.Client {
	cStorage, err := storage.NewClient(ctx)
	if err != nil {
		panic(err)
	}
	return cStorage
}

func NewBucketHandle(storageCli *storage.Client, bucketName string) *storage.BucketHandle {
	return storageCli.Bucket(bucketName)
}
