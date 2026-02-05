package repository

import (
	"context"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
)

//go:generate go tool go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type Asset interface {
	// GenerateWritePresignedURL selects bucket based on path prefix (public/ vs private/)
	GenerateWritePresignedURL(
		ctx context.Context,
		contentType model.ContentType,
		path string,
		expires time.Duration,
	) (string, error)

	// GenerateReadURL returns asset read URL
	// For private paths: returns presigned URL with rounded expiration for caching
	// For public paths: returns public base URL + path (no signing)
	GenerateReadURL(
		ctx context.Context,
		path string,
		requestTime time.Time,
	) (string, error)
}
