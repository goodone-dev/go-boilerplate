package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/cache"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
	htterror "github.com/goodone-dev/go-boilerplate/internal/utils/http_response/error"
)

type RateLimitMode string

const (
	SingleLimiter RateLimitMode = "single"
	GlobalLimiter RateLimitMode = "global"
)

func RateLimiterHandler(cache cache.Cache, limit int, ttl time.Duration, mode RateLimitMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		ctx, span := tracer.Start(c.Request.Context())
		defer func() {
			span.Stop(err)
		}()

		key := fmt.Sprintf("rate_limit:%s", c.ClientIP())
		if mode == SingleLimiter {
			key = fmt.Sprintf("%s:%s %s", key, c.Request.Method, c.Request.URL.Path)
		}

		val, err := cache.Get(ctx, key)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		if val.ToInt() >= limit {
			c.Error(htterror.NewTooManyRequestError("rate limit exceeded, please try again later"))
			c.Abort()
			return
		}

		_, err = cache.Incr(ctx, key)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		err = cache.Expire(ctx, key, ttl)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.Next()
	}
}
