package repository

import (
	"github.com/BagusAK95/go-skeleton/internal/domain/order"
	database "github.com/BagusAK95/go-skeleton/internal/infrastructure/database/sql"
)

type orderRepo struct {
	database.IBaseRepository[order.Order]
}

func NewOrderRepo(baseRepo database.IBaseRepository[order.Order]) order.IOrderRepository {
	return &orderRepo{
		baseRepo,
	}
}
