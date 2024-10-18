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
	requestTime time.Time,
) (*AssetCreatePresignedURLResult, error) {
	asset, err := model.NewAsset(
		assetType,
		contentType,
		requestTime,
	)
	if err != nil {
		return nil, err
	}
	presignedURL, err := s.assetRepository.GenerateWritePresignedURL(
		ctx,
		contentType,
		asset.Path,
		asset.Expiration(),
	)
	if err != nil {
		return nil, err
	}
	if err := s.assetPathCache.Set(
		ctx,
		asset,
	); err != nil {
		return nil, err
	}

	return &AssetCreatePresignedURLResult{
		AssetID:      asset.ID,
		PresignedURL: presignedURL,
	}, nil
}

func (s *assetService) GetWithValidate(
	ctx context.Context,
	assetType model.AssetType,
	assetID string,
) (string, error) {
	got, err := s.assetPathCache.Get(ctx, assetID)
	if err != nil {
		return "", err
	}
	if err := model.ValidateAssetPath(assetType, got); err != nil {
		return "", err
	}
	return got, nil
}

func (s *assetService) BatchSetStaffURLs(
	ctx context.Context,
	staffs model.Staffs,
) error {
	for _, staff := range staffs {
		imageURL, err := s.assetRepository.GenerateReadPresignedURL(
			ctx,
			staff.ImagePath,
			15*time.Minute,
		)
		if err != nil {
			return err
		}
		staff.SetImageURL(imageURL)
	}
	return nil
}
