package usecase

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/usecase/input"
	"github.com/abyssparanoia/rapid-go/internal/usecase/output"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_usecase
type StaffAssetInteractor interface {
	CreatePresignedURL(
		ctx context.Context,
		param *input.StaffCreateAssetPresignedURL,
	) (*output.StaffCreateAssetPresignedURL, error)
}
