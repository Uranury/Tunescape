package cache

import (
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, data []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
}

type cache struct {
	rds *redis.Client
}

func New(rds *redis.Client) Cache {
	return &cache{rds: rds}
}

func (c *cache) Get(ctx context.Context, key string) ([]byte, error) {
	return c.rds.Get(ctx, key).Bytes()
}

func (c *cache) Set(ctx context.Context, key string, data []byte, ttl time.Duration) error {
	return c.rds.Set(ctx, key, data, ttl).Err()
}

func (c *cache) Delete(ctx context.Context, key string) error {
	return c.rds.Del(ctx, key).Err()
}

func (c *cache) Exists(ctx context.Context, key string) (bool, error) {
	n, err := c.rds.Exists(ctx, key).Result()
	return n == 1, err
}
