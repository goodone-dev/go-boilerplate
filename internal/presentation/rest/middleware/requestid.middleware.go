package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

func RequestIdHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		span := trace.SpanFromContext(c.Request.Context())
		if span.SpanContext().HasTraceID() {
			c.Header("X-Request-Id", span.SpanContext().TraceID().String())
		} else {
			c.Header("X-Request-Id", uuid.NewString())
		}

		c.Next()
	}
}
