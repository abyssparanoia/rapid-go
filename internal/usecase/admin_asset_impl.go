package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/service"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/abyssparanoia/rapid-go/internal/usecase/output"
)

type adminAssetInteractor struct {
	assetService service.Asset
}

func NewAdminAssetInteractor(
	assetService service.Asset,
) AdminAssetInteractor {
	return &adminAssetInteractor{
		assetService: assetService,
	}
}

func (i *adminAssetInteractor) CreatePresignedURL(
	ctx context.Context,
	param *input.AdminCreateAssetPresignedURL,
) (*output.AdminCreateAssetPresignedURL, error) {
	if err := param.Validate(); err != nil {
		return nil, err
	}
	got, err := i.assetService.CreatePresignedURL(
		ctx,
		param.AssetType,
		param.ContentType,
		param.RequestTime,
	)
	if err != nil {
		return nil, err
	}
	return output.NewAdminCreateAssetPresignedURL(
		got.AssetKey,
		got.PresignedURL,
	), nil
}
