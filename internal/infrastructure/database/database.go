package database

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	MaxRetries     = 5
	InitialBackoff = 1 * time.Second
	MaxBackoff     = 30 * time.Second
)

type Entity interface {
	TableName() string
	RepositoryName() string
}

type BaseEntity[I any] struct {
	ID        I          `json:"id" bson:"_id,omitempty"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" bson:"deleted_at"`
}

func (b *BaseEntity[I]) BeforeCreate(tx *gorm.DB) (err error) {
	id, ok := any(b.ID).(uuid.UUID)
	if !ok {
		return nil
	} else if id != uuid.Nil {
		return nil
	}

	id, err = uuid.NewV7()
	if err != nil {
		return err
	}

	b.ID = any(id).(I)

	return nil
}

func RetryWithBackoff[C any](ctx context.Context, operation string, fn func() (C, error)) (res C, err error) {
	backoff := InitialBackoff

	for attempt := 0; attempt <= MaxRetries; attempt++ {
		if attempt > 0 {
			logger.Warnf(ctx, "⚠️ Retrying %s (attempt %d/%d) after %v", operation, attempt, MaxRetries, backoff)
			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return res, ctx.Err()
			}
		}

		res, err = fn()
		if err == nil {
			if attempt > 0 {
				logger.Infof(ctx, "✅ %s succeeded after %d attempts", operation, attempt+1)
			}
			return res, nil
		}

		if attempt < MaxRetries {
			backoff = min(time.Duration(float64(InitialBackoff)*math.Pow(2, float64(attempt))), MaxBackoff)
		}
	}

	return res, fmt.Errorf("%s failed after %d attempts: %w", operation, MaxRetries+1, err)
}
