package service

import (
	"context"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_service
type Asset interface {
	CreatePresignedURL(
		ctx context.Context,
		assetType model.AssetType,
		contentType model.ContentType,
		requestTime time.Time,
	) (*AssetCreatePresignedURLResult, error)
	GetWithValidate(
		ctx context.Context,
		assetType model.AssetType,
		assetID string,
	) (string, error)
	BatchSetStaffURLs(
		ctx context.Context,
		staffs model.Staffs,
	) error
}

type AssetCreatePresignedURLResult struct {
	AssetID      string
	PresignedURL string
}
