package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/utils/error"
)

func ContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), config.ContextTimeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		done := make(chan struct{})
		go func() {
			c.Next()
			close(done)
		}()

		select {
		case <-done:
			// Handler finished, do nothing
		case <-ctx.Done():
			if !c.Writer.Written() {
				switch ctx.Err() {
				case context.Canceled:
					c.Error(error.NewRequestTimeoutError("request canceled by user"))
				case context.DeadlineExceeded:
					c.Error(error.NewRequestTimeoutError("request timed out"))
				}
				c.Abort()
			}
		}
	}
}
