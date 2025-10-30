package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/cache"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	client *redis.Client
}

func createClient(ctx context.Context) (client *redis.Client) {
	options := &redis.Options{
		Addr: fmt.Sprintf("%v:%v", config.RedisConfig.Host, config.RedisConfig.Port),
		DB:   config.RedisConfig.DB,
	}

	if config.RedisConfig.TLS {
		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS13,
		}
	}

	if config.RedisConfig.Password != "" {
		options.Password = config.RedisConfig.Password
	}

	client = redis.NewClient(options)
	if _, err := client.Ping(ctx).Result(); err != nil {
		logger.Fatal(ctx, err, "❌ Failed to establish Redis connection")
	}

	return client
}

func NewClient(ctx context.Context) cache.ICache {
	return &redisClient{
		client: createClient(ctx),
	}
}

func (c *redisClient) Ping(ctx context.Context) (err error) {
	ctx, span := tracer.Start(ctx)
	defer func() {
		span.Stop(err)
	}()

	return c.client.Ping(ctx).Err()
}

func (c *redisClient) Get(ctx context.Context, key string) (res string, err error) {
	ctx, span := tracer.Start(ctx, key)
	defer func() {
		span.Stop(err, res)
	}()

	res, err = c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return
}

func (c *redisClient) Set(ctx context.Context, key string, val any, ttl time.Duration) (err error) {
	ctx, span := tracer.Start(ctx, key, val, ttl)
	defer func() {
		span.Stop(err)
	}()

	return c.client.Set(ctx, key, val, ttl).Err()
}

func (c *redisClient) TTL(ctx context.Context, key string) (ttl time.Duration, err error) {
	ctx, span := tracer.Start(ctx, key)
	defer func() {
		span.Stop(err, ttl)
	}()

	return c.client.TTL(ctx, key).Result()
}

func (c *redisClient) Del(ctx context.Context, keys ...string) (err error) {
	ctx, span := tracer.Start(ctx, keys)
	defer func() {
		span.Stop(err)
	}()

	return c.client.Del(ctx, keys...).Err()
}

func (c *redisClient) Incr(ctx context.Context, key string) (res int64, err error) {
	ctx, span := tracer.Start(ctx, key)
	defer func() {
		span.Stop(err, res)
	}()

	return c.client.Incr(ctx, key).Result()
}

func (c *redisClient) Decr(ctx context.Context, key string) (res int64, err error) {
	ctx, span := tracer.Start(ctx, key)
	defer func() {
		span.Stop(err, res)
	}()

	return c.client.Decr(ctx, key).Result()
}

func (c *redisClient) IncrBy(ctx context.Context, key string, value int64) (res int64, err error) {
	ctx, span := tracer.Start(ctx, key, value)
	defer func() {
		span.Stop(err, res)
	}()

	return c.client.IncrBy(ctx, key, value).Result()
}

func (c *redisClient) DecrBy(ctx context.Context, key string, value int64) (res int64, err error) {
	ctx, span := tracer.Start(ctx, key, value)
	defer func() {
		span.Stop(err, res)
	}()

	return c.client.DecrBy(ctx, key, value).Result()
}

func (c *redisClient) Shutdown(ctx context.Context) (err error) {
	return c.client.Close()
}
