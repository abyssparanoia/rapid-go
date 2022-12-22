package repository

import (
	"context"
	"time"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type Asset interface {
	GenerateWritePresignedURL(
		ctx context.Context,
		contentType string,
		path string,
		expires time.Duration,
	) (string, error)
	GenerateReadPresignedURL(
		ctx context.Context,
		path string,
		expires time.Duration,
	) (string, error)
}
