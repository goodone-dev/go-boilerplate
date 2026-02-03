package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/utils/http_response/error"
)

func ContextTimeoutHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), config.ContextTimeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		// Check if timeout occurred after handlers complete
		if ctx.Err() != nil && !c.Writer.Written() {
			switch ctx.Err() {
			case context.Canceled:
				c.Error(error.NewRequestTimeoutError("your request was canceled before completion"))
			case context.DeadlineExceeded:
				c.Error(error.NewRequestTimeoutError("request took too long to process, please try again"))
			}
			c.Abort()
		}
	}
}
