package service

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_service
type Asset interface {
	CreatePresignedURL(
		ctx context.Context,
		assetType model.AssetType,
		contentType string,
	) (*AssetCreatePresignedURLResult, error)
	GetWithValidate(
		ctx context.Context,
		assetType model.AssetType,
		assetKey string,
	) (string, error)
	BatchSetStaffURLs(
		ctx context.Context,
		staffs model.Staffs,
	) error
}

type AssetCreatePresignedURLResult struct {
	AssetKey     string
	PresignedURL string
}
