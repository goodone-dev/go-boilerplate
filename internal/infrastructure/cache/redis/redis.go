package redis

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"
	"time"

	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/cache"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
	"github.com/goodone-dev/go-boilerplate/internal/utils/retry"
	"github.com/redis/go-redis/extra/redisotel/v9"
	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	client *redis.Client
}

func createClient(ctx context.Context) (client *redis.Client) {
	redis.SetLogger(&noLogger{})

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

	_, err := retry.RetryWithBackoff(ctx, "Redis connection test", func() (any, error) {
		return nil, client.Ping(ctx).Err()
	})
	if err != nil {
		logger.Fatal(ctx, err, "❌ Redis failed to establish connection").Write()
	}

	traceOpts := redisotel.WithCommandFilter(func(cmd redis.Cmder) bool {
		return !strings.EqualFold(cmd.Name(), "ping") // Skip tracing PING commands
	})

	if err := redisotel.InstrumentTracing(client, traceOpts); err != nil {
		logger.Fatal(ctx, err, "❌ Redis failed to instrument connection").Write()
	}

	return client
}

func NewClient(ctx context.Context) cache.Cache {
	client := &redisClient{
		client: createClient(ctx),
	}

	go client.Monitor(ctx)

	return client
}

func (c *redisClient) Ping(ctx context.Context) (err error) {
	return c.client.Ping(ctx).Err()
}

func (c *redisClient) Get(ctx context.Context, key string) (res *cache.CacheValue, err error) {
	ctx, span := tracer.Start(ctx, key)
	defer func() {
		span.Stop(err, res)
	}()

	str, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	val := cache.CacheValue(str)
	return &val, nil
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

func (c *redisClient) Expire(ctx context.Context, key string, ttl time.Duration) (err error) {
	ctx, span := tracer.Start(ctx, key, ttl)
	defer func() {
		span.Stop(err)
	}()

	return c.client.ExpireNX(ctx, key, ttl).Err()
}

func (c *redisClient) Shutdown(ctx context.Context) (err error) {
	return c.client.Close()
}

func (c *redisClient) Monitor(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	var wasLost bool

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := c.Ping(ctx)
			if err != nil {
				if !wasLost {
					logger.Errorf(ctx, err, "🛑 Redis connection lost").Write()
					wasLost = true
				}
			} else {
				if wasLost {
					logger.Info(ctx, "✅ Redis connection restored").Write()
					wasLost = false
				}
			}
		}
	}
}

// noLogger is a no-op logger that implements redis internal.Logging interface
type noLogger struct{}

func (n *noLogger) Printf(_ context.Context, _ string, _ ...any) {}
