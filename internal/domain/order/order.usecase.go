package order

import (
	"context"
)

type OrderUsecase interface {
	Create(ctx context.Context, req CreateOrderRequest) (*CreateOrderResponse, error)
}
