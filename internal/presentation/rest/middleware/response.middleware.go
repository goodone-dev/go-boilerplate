package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func ResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())
		if span.SpanContext().HasTraceID() {
			c.Header("X-Request-Id", span.SpanContext().TraceID().String())
		}

		c.Next()
	}
}
