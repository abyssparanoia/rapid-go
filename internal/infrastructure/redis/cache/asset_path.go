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
	authContext model.AssetAuthContext,
) string {
	return fmt.Sprintf("asset_path:%s:%s", assetKey, authContext.String())
}

func (c *assetPath) Get(
	ctx context.Context,
	id string,
	authContext model.AssetAuthContext,
) (string, error) {
	cacheKey := c.buildCacheKey(id, authContext)
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
	cacheKey := c.buildCacheKey(asset.ID, asset.AuthContext)
	if err := c.cli.Set(ctx, cacheKey, asset.Path, asset.Expiration()).Err(); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}

func (c *assetPath) Clear(
	ctx context.Context,
	id string,
) error {
	// Delete all cache entries for this asset ID regardless of auth context
	pattern := fmt.Sprintf("asset_path:%s:*", id)
	iter := c.cli.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := c.cli.Del(ctx, iter.Val()).Err(); err != nil {
			return errors.InternalErr.Wrap(err)
		}
	}
	if err := iter.Err(); err != nil {
		return errors.InternalErr.Wrap(err)
	}
	return nil
}
