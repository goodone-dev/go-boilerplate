package order

import (
	"context"
)

type IOrderUsecase interface {
	Create(ctx context.Context, req CreateOrderRequest) (*Order, error)
}
