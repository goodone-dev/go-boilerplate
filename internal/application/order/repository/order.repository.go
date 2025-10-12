package repository

import (
	"github.com/goodonedev/go-boilerplate/internal/domain/order"
	"github.com/goodonedev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository struct {
	database.IBaseRepository[gorm.DB, uuid.UUID, order.Order]
}

func NewOrderRepo(baseRepo database.IBaseRepository[gorm.DB, uuid.UUID, order.Order]) order.IOrderRepository {
	return &OrderRepository{
		baseRepo,
	}
}
