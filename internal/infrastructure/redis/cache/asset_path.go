package cache

import (
	"context"
	"fmt"

	"github.com/abyssparanoia/rapid-go/internal/domain/cache"
	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/abyssparanoia/rapid-go/internal/domain/model"
	"github.com/redis/go-redis/v9"
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

func (c *assetPath) Get(
	ctx context.Context,
	id string,
) (string, error) {
	cacheKey := c.buildCacheKey(id)
	got, err := c.cli.Get(ctx, cacheKey).Result()
	if err != nil {
		if err == redis.Nil {
			return "", errors.AssetNotFoundErr.Wrap(err)
		}
		return "", errors.InternalErr.Wrap(err)
	}
	return got, nil
}

func (c *assetPath) Set(
	ctx context.Context,
	asset *model.Asset,
) error {
	cacheKey := c.buildCacheKey(asset.ID)
	if err := c.cli.Set(ctx, cacheKey, asset.Path, asset.Expiration()).Err(); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}

func (c *assetPath) Clear(
	ctx context.Context,
	id string,
) error {
	cacheKey := c.buildCacheKey(id)
	if err := c.cli.Del(ctx, cacheKey).Err(); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}
