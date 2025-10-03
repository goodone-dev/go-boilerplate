package middleware

import (
	"fmt"
	"strconv"
	"time"

	"github.com/BagusAK95/go-boilerplate/internal/infrastructure/cache"
	httError "github.com/BagusAK95/go-boilerplate/internal/utils/error"
	"github.com/BagusAK95/go-boilerplate/internal/utils/tracer"
	"github.com/gin-gonic/gin"
)

func RateLimitMiddleware(cache cache.ICache, limit int, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		ctx, span := tracer.StartSpan(c.Request.Context())
		defer func() {
			span.EndSpan(err)
		}()

		key := fmt.Sprintf("rate_limit:%s", c.ClientIP())

		countStr, err := cache.Get(ctx, key)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		if countStr == "" {
			if err := cache.Set(ctx, key, 1, duration); err != nil {
				c.Error(err)
				c.Abort()
				return
			}

			c.Next()
			return
		}

		countInt, _ := strconv.Atoi(countStr)
		if countInt >= limit {
			c.Error(httError.NewTooManyRequestError("too many request"))
			c.Abort()
			return
		}

		_, err = cache.Incr(ctx, key)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		c.Next()
	}
}
