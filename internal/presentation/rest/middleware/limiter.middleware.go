package middleware

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/cache"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
	htterror "github.com/goodone-dev/go-boilerplate/internal/utils/http_response/error"
)

func SingleLimiterMiddleware(cache cache.ICache, limit int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		ctx, span := tracer.StartSpan(c.Request.Context())
		defer func() {
			span.EndSpan(err)
		}()

		key := fmt.Sprintf("rate_limit:%s:%s %s", c.ClientIP(), c.Request.Method, c.Request.URL.Path)

		err = handleLimiter(ctx, cache, key, limit, duration)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.Next()
	}
}

func GlobalLimiterMiddleware(cache cache.ICache, limit int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		ctx, span := tracer.StartSpan(c.Request.Context())
		defer func() {
			span.EndSpan(err)
		}()

		key := fmt.Sprintf("rate_limit:%s", c.ClientIP())

		err = handleLimiter(ctx, cache, key, limit, duration)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.Next()
	}
}

func handleLimiter(ctx context.Context, cache cache.ICache, key string, limit int, duration time.Duration) error {
	countStr, err := cache.Get(ctx, key)
	if err != nil {
		return err
	}

	if countStr == "" {
		if err := cache.Set(ctx, key, 1, duration); err != nil {
			return err
		}

		return nil
	}

	countInt, _ := strconv.Atoi(countStr)
	if countInt >= limit {
		return htterror.NewTooManyRequestError("rate limit exceeded, please try again later")
	}

	_, err = cache.Incr(ctx, key)
	if err != nil {
		return err
	}

	return nil
}
