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
	contentType model.ContentType,
	authContext model.AssetAuthContext,
	requestTime time.Time,
) (*AssetCreatePresignedURLResult, error) {
	asset := model.NewAsset(
		assetType,
		contentType,
		authContext,
		requestTime,
	)
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
	authContext model.AssetAuthContext,
) (string, error) {
	got, err := s.assetPathCache.Get(ctx, assetID, authContext)
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
	requestTime time.Time,
) error {
	for _, staff := range staffs {
		imageURL, err := s.assetRepository.GenerateReadURL(
			ctx,
			staff.ImagePath,
			requestTime,
		)
		if err != nil {
			return err
		}
		staff.SetImageURL(imageURL)

		if staff.ReadonlyReference != nil && staff.ReadonlyReference.Tenant != nil {
			if err := s.BatchSetTenantURLs(ctx, model.Tenants{staff.ReadonlyReference.Tenant}, requestTime); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *assetService) BatchSetTenantURLs(
	ctx context.Context,
	tenants model.Tenants,
	requestTime time.Time,
) error {
	for _, tenant := range tenants {
		// no URLs to set for tenant yet
		if err := s.BatchSetTenantTagURLs(ctx, tenant.Tags, requestTime); err != nil {
			return err
		}
	}
	return nil
}

func (s *assetService) BatchSetTenantTagURLs(
	ctx context.Context,
	tenantTags model.TenantTags,
	requestTime time.Time,
) error {
	for range tenantTags {
		// no URLs to set for tenant tag yet
	}
	return nil
}
