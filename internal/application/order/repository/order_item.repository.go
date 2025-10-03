package repository

import (
	"github.com/BagusAK95/go-skeleton/internal/domain/order"
	database "github.com/BagusAK95/go-skeleton/internal/infrastructure/database/sql"
)

type orderItemRepo struct {
	database.IBaseRepository[order.OrderItem]
}

func NewOrderItemRepo(baseRepo database.IBaseRepository[order.OrderItem]) order.IOrderItemRepository {
	return &orderItemRepo{
		baseRepo,
	}
}
