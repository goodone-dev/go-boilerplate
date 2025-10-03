package repository

import (
	"github.com/BagusAK95/go-boilerplate/internal/domain/order"
	database "github.com/BagusAK95/go-boilerplate/internal/infrastructure/database/sql"
)

type OrderItemRepository struct {
	database.IBaseRepository[order.OrderItem]
}

func NewOrderItemRepo(baseRepo database.IBaseRepository[order.OrderItem]) order.IOrderItemRepository {
	return &OrderItemRepository{
		baseRepo,
	}
}
