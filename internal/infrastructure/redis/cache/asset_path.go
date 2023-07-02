package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/abyssparanoia/rapid-go/internal/domain/cache"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/go-redis/redis/v8"
)

type assetPath struct {
	cli *redis.Client
}

func NewAssetPath(
	cli *redis.Client,
) cache.AssetPath {
	return &assetPath{
		cli: cli,
	}
}

func (c *assetPath) buildCacheKey(
	assetKey string,
) string {
	return fmt.Sprintf("asset_path:%s", assetKey)
}

func (c *assetPath) GetWithValidate(
	ctx context.Context,
	assetType model.AssetType,
	assetKey string,
) (string, error) {
	cacheKey := c.buildCacheKey(assetKey)
	got, err := c.cli.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.AssetNotFoundErr.Wrap(err)
		}
		return "", errors.InternalErr.Wrap(err)
	}
	if err := model.ValidateAssetPath(assetType, got); err != nil {
		return "", err
	}
	return got, nil
}

func (c *assetPath) Set(
	ctx context.Context,
	asset *model.Asset,
	expiration time.Duration,
) error {
	cacheKey := c.buildCacheKey(asset.Key)
	if err := c.cli.Set(ctx, cacheKey, asset.Path, expiration).Err(); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}
