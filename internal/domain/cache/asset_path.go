package cache

import (
	"context"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
)

//go:generate go run go.uber.org/mock/mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_cache
type AssetPath interface {
	Get(
		ctx context.Context,
		id string,
	) (string, error)
	Set(
		ctx context.Context,
		asset *model.Asset,
	) error
	Clear(
		ctx context.Context,
		id string,
	) error
}
