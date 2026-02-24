package gcs

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func NewClient(ctx context.Context, emulatorHost string) *storage.Client {
	var opts []option.ClientOption

	if emulatorHost != "" {
		// Emulator mode: disable authentication and set custom endpoint
		opts = append(opts,
			option.WithoutAuthentication(),
			option.WithEndpoint(fmt.Sprintf("%s/storage/v1/", emulatorHost)),
		)
	}

	cStorage, err := storage.NewClient(ctx, opts...)
	if err != nil {
		panic(err)
	}
	return cStorage
}

func NewBucketHandle(storageCli *storage.Client, bucketName string) *storage.BucketHandle {
	return storageCli.Bucket(bucketName)
}
