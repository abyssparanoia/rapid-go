package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/abyssparanoia/rapid-go/internal/usecase/output"
)

type AdminAssetInteractor interface {
	CreatePresignedURL(
		ctx context.Context,
		param *input.AdminCreateAssetPresignedURL,
	) (*output.AdminCreateAssetPresignedURL, error)
}
