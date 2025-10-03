package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"time"

	"github.com/BagusAK95/go-boilerplate/internal/config"
	"github.com/BagusAK95/go-boilerplate/internal/infrastructure/cache"
	"github.com/BagusAK95/go-boilerplate/internal/utils/tracer"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func createClient(ctx context.Context) (client *redis.Client) {
	options := &redis.Options{
		Addr: fmt.Sprintf("%v:%v", config.RedisConfig.Host, config.RedisConfig.Port),
		DB:   config.RedisConfig.DB,
	}

	if config.RedisConfig.TLS {
		options.TLSConfig = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	if config.RedisConfig.Password != "" {
		options.Password = config.RedisConfig.Password
	}

	client = redis.NewClient(options)
	if _, err := client.Ping(ctx).Result(); err != nil {
		log.Fatalf("❌ Could not connect to Redis: %v", err)
	}

	return client
}

func NewClient(ctx context.Context) cache.ICache {
	return &RedisClient{
		client: createClient(ctx),
	}
}

func (c *RedisClient) Ping(ctx context.Context) (err error) {
	ctx, span := tracer.StartSpan(ctx)
	defer func() {
		span.EndSpan(err)
	}()

	return c.client.Ping(ctx).Err()
}

func (c *RedisClient) Get(ctx context.Context, key string) (res string, err error) {
	ctx, span := tracer.StartSpan(ctx, key)
	defer func() {
		span.EndSpan(err, res)
	}()

	res, err = c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	} else if err != nil {
		return "", err
	}

	return
}

func (c *RedisClient) Set(ctx context.Context, key string, val any, ttl time.Duration) (err error) {
	ctx, span := tracer.StartSpan(ctx, key, val, ttl)
	defer func() {
		span.EndSpan(err)
	}()

	return c.client.Set(ctx, key, val, ttl).Err()
}

func (c *RedisClient) TTL(ctx context.Context, key string) (ttl time.Duration, err error) {
	ctx, span := tracer.StartSpan(ctx, key)
	defer func() {
		span.EndSpan(err, ttl)
	}()

	return c.client.TTL(ctx, key).Result()
}

func (c *RedisClient) Del(ctx context.Context, keys ...string) (err error) {
	ctx, span := tracer.StartSpan(ctx, keys)
	defer func() {
		span.EndSpan(err)
	}()

	return c.client.Del(ctx, keys...).Err()
}

func (c *RedisClient) Incr(ctx context.Context, key string) (res int64, err error) {
	ctx, span := tracer.StartSpan(ctx, key)
	defer func() {
		span.EndSpan(err, res)
	}()

	return c.client.Incr(ctx, key).Result()
}

func (c *RedisClient) Decr(ctx context.Context, key string) (res int64, err error) {
	ctx, span := tracer.StartSpan(ctx, key)
	defer func() {
		span.EndSpan(err, res)
	}()

	return c.client.Decr(ctx, key).Result()
}

func (c *RedisClient) IncrBy(ctx context.Context, key string, value int64) (res int64, err error) {
	ctx, span := tracer.StartSpan(ctx, key, value)
	defer func() {
		span.EndSpan(err, res)
	}()

	return c.client.IncrBy(ctx, key, value).Result()
}

func (c *RedisClient) DecrBy(ctx context.Context, key string, value int64) (res int64, err error) {
	ctx, span := tracer.StartSpan(ctx, key, value)
	defer func() {
		span.EndSpan(err, res)
	}()

	return c.client.DecrBy(ctx, key, value).Result()
}

func (c *RedisClient) Close() (err error) {
	return c.client.Close()
}
