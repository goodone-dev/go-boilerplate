package order

import (
	"context"
)

type IOrderUsecase interface {
	Create(ctx context.Context, req CreateOrderRequest) (*CreateOrderResponse, error)
}
