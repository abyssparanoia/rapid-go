package redis_helper

import (
	"context"
	"fmt"

	"github.com/abyssparanoia/rapid-go/internal/domain/errors"
	"github.com/cenkalti/backoff"
	"github.com/redis/go-redis/v9"
)

func RunTransaction(ctx context.Context, redisClient *redis.Client, fn func(ctx context.Context, tx *redis.Tx) error, keys ...string) error {
	b := backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 100)

	if err := backoff.Retry(func() error {
		err := redisClient.Watch(ctx, func(tx *redis.Tx) error {
			return fn(ctx, tx)
		}, keys...)
		if err == redis.TxFailedErr {
			return errors.InternalErr.Wrap(err)
		}
		if err != nil {
			return backoff.Permanent(err)
		}
		return nil
	}, b); err != nil {
		return fmt.Errorf("transaction failed: %w", err)
	}
	return nil
}
