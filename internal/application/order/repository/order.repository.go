package repository

import (
	"github.com/BagusAK95/go-skeleton/internal/domain/order"
	database "github.com/BagusAK95/go-skeleton/internal/infrastructure/database/sql"
)

type OrderRepository struct {
	database.IBaseRepository[order.Order]
}

func NewOrderRepo(baseRepo database.IBaseRepository[order.Order]) order.IOrderRepository {
	return &OrderRepository{
		baseRepo,
	}
}
