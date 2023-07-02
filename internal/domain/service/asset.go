package service

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_service
type Asset interface {
	CreatePresignedURL(
		ctx context.Context,
		assetType model.AssetType,
		contentType string,
	) (*AssetCreatePresignedURLResult, error)
}

type AssetCreatePresignedURLResult struct {
	AssetKey     string
	PresignedURL string
}
