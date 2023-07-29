package cache

import (
	"context"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/model"
)

//go:generate mockgen -source=$GOFILE -destination=mock/$GOFILE -package=mock_cache
type AssetPath interface {
	Get(
		ctx context.Context,
		assetKey string,
	) (string, error)
	Set(
		ctx context.Context,
		asset *model.Asset,
		expiration time.Duration,
	) error
}
