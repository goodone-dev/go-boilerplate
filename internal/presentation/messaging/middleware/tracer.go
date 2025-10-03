package middleware

import (
	"context"
	"fmt"

	"github.com/BagusAK95/go-skeleton/internal/utils/tracer"
)

func TracerMiddleware[T any](topic string, handler func(context.Context, T) error) func(T) {
	return func(msg T) {
		ctx, span := tracer.SpanCustomName(fmt.Sprintf("MESSAGE %s", topic)).StartSpan(context.Background())

		err := handler(ctx, msg)
		defer func() {
			span.EndSpan(err)
		}()
	}
}
