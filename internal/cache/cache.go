package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"hr-system/internal/common/errors"
)

type Cache struct {
	rdb *redis.Client
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		rdb: client,
	}
}

func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func (c *Cache) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	err := c.rdb.Set(ctx, key, value, expiration).Err()
	return err
}

func (c *Cache) Del(ctx context.Context, key string) error {
	err := c.rdb.Del(ctx, key).Err()
	return err
}

func (c *Cache) DelByPrefix(ctx context.Context, prefix string) error {
	var cursor uint64
	var keys []string

	for {
		var scanKeys []string
		var err error
		scanKeys, cursor, err = c.rdb.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return fmt.Errorf("failed to scan keys: %w", err)
		}

		keys = append(keys, scanKeys...)

		if cursor == 0 {
			break
		}
	}

	errs := make([]error, 0, len(keys))
	if len(keys) > 0 {
		_, err := c.rdb.Del(ctx, keys...).Result()
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to delete keys: %w", err))
		}
	}

	return errors.Combine(errs...)
}
