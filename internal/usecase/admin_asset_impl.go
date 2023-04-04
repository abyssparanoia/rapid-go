package usecase

import (
	"context"
	"fmt"
	"mime"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/repository"
	"github.com/abyssparanoia/rapid-go/internal/pkg/uuid"
	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/abyssparanoia/rapid-go/internal/usecase/output"
)

type adminAssetInteractor struct {
	assetRepository repository.Asset
}

func NewAdminAssetInteractor(
	assetRepository repository.Asset,
) AdminAssetInteractor {
	return &adminAssetInteractor{
		assetRepository,
	}
}

func (i *adminAssetInteractor) CreatePresignedURL(
	ctx context.Context,
	param *input.AdminCreateAssetPresignedURL,
) (*output.AdminCreateAssetPresignedURL, error) {
	ext, err := mime.ExtensionsByType(param.ContentType)
	if err != nil {
		return nil, errors.InternalErr.Wrap(err)
	}
	path := fmt.Sprintf("%s/%s%s", param.AssetType.String(), uuid.UUIDBase64(), ext[0])
	presignedURL, err := i.assetRepository.GenerateWritePresignedURL(
		ctx,
		param.ContentType,
		path,
		// 有効期限は15分とする
		15*time.Minute,
	)
	if err != nil {
		return nil, err
	}
	return output.NewAdminCreateAssetPresignedURL(
		path,
		presignedURL,
	), nil
}
