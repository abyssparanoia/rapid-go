package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
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
	authContext := model.NewAdminAssetAuthContext(param.AdminID)
	got, err := i.assetService.CreatePresignedURL(
		ctx,
		param.AssetType,
		param.ContentType,
		authContext,
		param.RequestTime,
	)
	if err != nil {
		return nil, err
	}
	return output.NewAdminCreateAssetPresignedURL(
		got.AssetID,
		got.PresignedURL,
	), nil
}
