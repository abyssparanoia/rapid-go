package repository

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
)

const (
	// urlRoundingDuration is the time rounding unit for presigned URL cache optimization
	urlRoundingDuration = 5 * time.Minute
)

type asset struct {
	privateBucketHandle *storage.BucketHandle
	publicBucketHandle  *storage.BucketHandle
	publicAssetBaseURL  string
}

func NewAsset(
	privateBucketHandle *storage.BucketHandle,
	publicBucketHandle *storage.BucketHandle,
	publicAssetBaseURL string,
) repository.Asset {
	return &asset{
		privateBucketHandle: privateBucketHandle,
		publicBucketHandle:  publicBucketHandle,
		publicAssetBaseURL:  strings.TrimSuffix(publicAssetBaseURL, "/"),
	}
}

func (r *asset) GenerateWritePresignedURL(
	ctx context.Context,
	contentType model.ContentType,
	path string,
	expires time.Duration,
) (string, error) {
	// Select bucket based on path prefix
	bucketHandle := r.privateBucketHandle
	if strings.HasPrefix(path, "public/") {
		bucketHandle = r.publicBucketHandle
	}

	opts := &storage.SignedURLOptions{
		Expires:     time.Now().Add(expires),
		Method:      http.MethodPut,
		ContentType: contentType.String(),
	}
	signedURL, err := bucketHandle.SignedURL(path, opts)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	return signedURL, nil
}

func (r *asset) GenerateReadURL(
	ctx context.Context,
	path string,
	requestTime time.Time,
) (string, error) {
	// Public: return base URL + path (no signing)
	if strings.HasPrefix(path, "public/") {
		return fmt.Sprintf("%s/%s", r.publicAssetBaseURL, path), nil
	}

	// Private: generate presigned URL with rounded expiration for cache optimization
	// Round request time to 5-minute intervals so same URL is generated within window
	roundedTime := requestTime.Truncate(urlRoundingDuration)
	// Expiration: rounded time + 2 * rounding duration (ensures URL valid for full window)
	expiresAt := roundedTime.Add(urlRoundingDuration * 2)

	opts := &storage.SignedURLOptions{
		Expires: expiresAt,
		Method:  http.MethodGet,
	}
	signedURL, err := r.privateBucketHandle.SignedURL(path, opts)
	if err != nil {
		return "", errors.InternalErr.Wrap(err)
	}
	return signedURL, nil
}
