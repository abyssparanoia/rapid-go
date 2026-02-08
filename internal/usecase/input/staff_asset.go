package input

import (
	"fmt"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/abyssparanoia/rapid-go/internal/pkg/validation"
)

type StaffCreateAssetPresignedURL struct {
	StaffID     string            `validate:"required"`
	ContentType model.ContentType `validate:"required"`
	AssetType   model.AssetType   `validate:"required"`
	RequestTime time.Time         `validate:"required"`
}

func NewStaffCreateAssetPresignedURL(
	staffID string,
	contentType model.ContentType,
	assetType model.AssetType,
	requestTime time.Time,
) *StaffCreateAssetPresignedURL {
	return &StaffCreateAssetPresignedURL{
		StaffID:     staffID,
		ContentType: contentType,
		AssetType:   assetType,
		RequestTime: requestTime,
	}
}

func (p *StaffCreateAssetPresignedURL) Validate() error {
	if err := validation.Validate(p); err != nil {
		return errors.RequestInvalidArgumentErr.Wrap(err)
	}
	if !p.AssetType.Valid() {
		return errors.RequestInvalidArgumentErr.New().
			WithDetail(fmt.Sprintf("invalid asset type %s", p.AssetType))
	}
	if !p.ContentType.Valid() {
		return errors.RequestInvalidArgumentErr.New().
			WithDetail(fmt.Sprintf("invalid content type %s", p.ContentType))
	}
	return nil
}
