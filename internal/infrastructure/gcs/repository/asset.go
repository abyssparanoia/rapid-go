package repository

import (
	"context"
	"net/http"
	"time"

	"cloud.google.com/go/storage"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/pkg/now"
)

type asset struct {
	bucketHandle *storage.BucketHandle
}

func NewAsset(
	bucketHandle *storage.BucketHandle,
) repository.Asset {
	return &asset{
		bucketHandle,
	}
}

func (r *asset) GenerateWritePresignedURL(
	ctx context.Context,
	contentType string,
	path string,
	expires time.Duration,
) (string, error) {
	now := now.Now()
	opts := &storage.SignedURLOptions{
		Expires:     now.Add(expires),
		Method:      http.MethodPut,
		ContentType: contentType,
	}
	singedURL, err := r.bucketHandle.SignedURL(path, opts)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	return singedURL, nil
}

func (r *asset) GenerateReadPresignedURL(
	ctx context.Context,
	path string,
	expires time.Duration,
) (string, error) {
	now := now.Now()
	opts := &storage.SignedURLOptions{
		Expires: now.Add(expires),
		Method:  http.MethodGet,
	}
	singedURL, err := r.bucketHandle.SignedURL(path, opts)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	return singedURL, nil
}
