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
	backoff := config.RetryBackoff.InitialBackoff

	for attempt := 0; attempt <= config.RetryBackoff.MaxRetries; attempt++ {
		if attempt > 0 {
			logger.Warnf(ctx, "🔁 Retrying %s (attempt %d/%d) after %v", operation, attempt, config.RetryBackoff.MaxRetries, backoff).Write()
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return res, ctx.Err()
			}
		}

		res, err = fn()
		if err == nil {
			if attempt > 0 {
				logger.Infof(ctx, "✅ %s succeeded after %d attempts", operation, attempt+1).Write()
			}
			return res, nil
		}

		if attempt < config.RetryBackoff.MaxRetries {
			backoff = min(time.Duration(float64(config.RetryBackoff.InitialBackoff)*math.Pow(2, float64(attempt))), config.RetryBackoff.MaxBackoff)
		}
	}

	return res, fmt.Errorf("%s failed after %d attempts: %w", operation, config.RetryBackoff.MaxRetries+1, err)
}
