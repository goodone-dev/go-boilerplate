package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/cache"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
)

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseWriter) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

func IdempotencyMiddleware(cache cache.ICache, duration time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		ctx, span := tracer.StartSpan(c.Request.Context())
		defer func() {
			span.EndSpan(err)
		}()

		idempotencyKey := c.GetHeader("X-Idempotency-Key")
		if idempotencyKey == "" {
			c.Next()
			return
		}

		key := "idempotency:" + idempotencyKey

		bodyStr, err := cache.Get(ctx, key)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		if bodyStr != "" {
			var bodyObj any
			if err = json.Unmarshal([]byte(bodyStr), &bodyObj); err == nil {
				c.JSON(http.StatusOK, bodyObj)
				c.Abort()
				return
			}
		}

		blw := &responseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}

		c.Writer = blw
		c.Next()

		if blw.body.String() == "" {
			c.Abort()
			return
		}

		err = cache.Set(ctx, key, blw.body.String(), duration)
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}

		_, err = blw.ResponseWriter.Write(blw.body.Bytes())
		if err != nil {
			c.Error(err)
			c.Abort()
			return
		}
	}
}
