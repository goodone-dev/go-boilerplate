package retry

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
)

func RetryWithBackoff[D any](ctx context.Context, operation string, fn func() (D, error)) (res D, err error) {
	backoff := config.RetryConfig.InitialBackoff

	for attempt := 0; attempt <= config.RetryConfig.MaxRetries; attempt++ {
		if attempt > 0 {
			logger.With().Warnf(ctx, "🔁 Retrying %s (attempt %d/%d) after %v", operation, attempt, config.RetryConfig.MaxRetries, backoff)
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return res, ctx.Err()
			}
		}

		res, err = fn()
		if err == nil {
			if attempt > 0 {
				logger.With().Infof(ctx, "✅ %s succeeded after %d attempts", operation, attempt+1)
			}
			return res, nil
		}

		if attempt < config.RetryConfig.MaxRetries {
			backoff = min(time.Duration(float64(config.RetryConfig.InitialBackoff)*math.Pow(2, float64(attempt))), config.RetryConfig.MaxBackoff)
		}
	}

	return res, fmt.Errorf("%s failed after %d attempts: %w", operation, config.RetryConfig.MaxRetries+1, err)
}
