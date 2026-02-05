package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/service"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/abyssparanoia/rapid-go/internal/usecase/output"
)

type staffAssetInteractor struct {
	assetService service.Asset
}

func NewStaffAssetInteractor(
	assetService service.Asset,
) StaffAssetInteractor {
	return &staffAssetInteractor{
		assetService: assetService,
	}
}

func (i *staffAssetInteractor) CreatePresignedURL(
	ctx context.Context,
	param *input.StaffCreateAssetPresignedURL,
) (*output.StaffCreateAssetPresignedURL, error) {
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
	return output.NewStaffCreateAssetPresignedURL(
		got.AssetID,
		got.PresignedURL,
	), nil
}
