package order

import (
	"context"

	"github.com/google/uuid"
)

type OrderUsecase interface {
	Create(ctx context.Context, req CreateOrderRequest) (*CreateOrderResponse, error)
	GetDetail(ctx context.Context, id uuid.UUID) (*DetailOrderResponse, error)
}
