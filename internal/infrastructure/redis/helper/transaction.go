package redis_helper

import (
	"context"

	"github.com/cenkalti/backoff"
	"github.com/redis/go-redis/v9"
)

func RunTransaction(ctx context.Context, redisClient *redis.Client, fn func(ctx context.Context, tx *redis.Tx) error, keys ...string) error {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 100)

	return backoff.Retry(func() error {
		err := redisClient.Watch(ctx, func(tx *redis.Tx) error {
			return fn(ctx, tx)
		}, keys...)
		if err == redis.TxFailedErr {
			return err
		}
		if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}, b)
}
