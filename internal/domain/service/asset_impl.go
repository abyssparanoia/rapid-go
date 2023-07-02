package service

import (
	"context"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/cache"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
)

type assetService struct {
	assetRepository repository.Asset
	assetPathCache  cache.AssetPath
}

func NewAsset(
	assetRepository repository.Asset,
	assetPathCache cache.AssetPath,
) Asset {
	return &assetService{
		assetRepository: assetRepository,
		assetPathCache:  assetPathCache,
	}
}

func (s *assetService) CreatePresignedURL(
	ctx context.Context,
	assetType model.AssetType,
	contentType string,
) (*AssetCreatePresignedURLResult, error) {
	asset, err := model.NewAsset(
		assetType,
		contentType,
	)
	if err != nil {
		return nil, err
	}
	presignedURL, err := s.assetRepository.GenerateWritePresignedURL(
		ctx,
		contentType,
		asset.Path,
		// 有効期限は15分とする
		15*time.Minute,
	)
	if err != nil {
		return nil, err
	}
	if err := s.assetPathCache.Set(
		ctx,
		asset,
		// 有効期限は30分とする
		24*time.Hour,
	); err != nil {
		return nil, err
	}

	return &AssetCreatePresignedURLResult{
		AssetKey:     asset.Key,
		PresignedURL: presignedURL,
	}, nil
}
