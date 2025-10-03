package repository

import (
	"github.com/BagusAK95/go-boilerplate/internal/domain/order"
	database "github.com/BagusAK95/go-boilerplate/internal/infrastructure/database/sql"
)

type OrderRepository struct {
	database.IBaseRepository[order.Order]
}

func NewOrderRepo(baseRepo database.IBaseRepository[order.Order]) order.IOrderRepository {
	return &OrderRepository{
		baseRepo,
	}
}
