package cache

import (
	"context"
	"time"
)

type Cache interface {
	Ping(ctx context.Context) error
	Get(ctx context.Context, key string) (*CacheValue, error)
	Set(ctx context.Context, key string, val any, ttl time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
	Del(ctx context.Context, keys ...string) error
	Incr(ctx context.Context, key string) (int64, error)
	Decr(ctx context.Context, key string) (int64, error)
	IncrBy(ctx context.Context, key string, value int64) (int64, error)
	DecrBy(ctx context.Context, key string, value int64) (int64, error)
	Expire(ctx context.Context, key string, ttl time.Duration) error
	Shutdown(ctx context.Context) error
}
