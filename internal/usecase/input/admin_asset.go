package input

import (
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/pkg/validation"
)

type AdminCreateAssetPresignedURL struct {
	ContentType string          `validate:"required"`
	AssetType   model.AssetType `validate:"required"`
}

func NewAdminCreateAssetPresignedURL(
	contentType string,
	assetType model.AssetType,
) *AdminCreateAssetPresignedURL {
	return &AdminCreateAssetPresignedURL{
		ContentType: contentType,
		AssetType:   assetType,
	}
}

func (p *AdminCreateAssetPresignedURL) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	if !p.AssetType.Valid() {
		return errors.RequestInvalidArgumentErr.Errorf("invalid asset type %s", p.AssetType)
	}
	return nil
}
