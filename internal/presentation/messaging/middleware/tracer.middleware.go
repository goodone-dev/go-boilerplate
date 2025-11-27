package middleware

import (
	"context"
	"fmt"

	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
)

func TracerMiddleware[T any](topic string, handler func(context.Context, T) error) func(T) {
	return func(msg T) {
		ctx, span := tracer.CustomName(fmt.Sprintf("MESSAGE %s", topic)).Start(context.Background())

		err := handler(ctx, msg)
		defer func() {
			span.Stop(err)
		}()
	}
}
