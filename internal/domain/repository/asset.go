package repository

import (
	"context"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_repository
type Asset interface {
	GenerateWritePresignedURL(
		ctx context.Context,
		contentType model.ContentType,
		path string,
		expires time.Duration,
	) (string, error)
	GenerateReadPresignedURL(
		ctx context.Context,
		path string,
		expires time.Duration,
	) (string, error)
}
