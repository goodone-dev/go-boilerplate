package repository

import (
	"github.com/goodone-dev/go-boilerplate/internal/domain/order"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderRepository struct {
	database.IBaseRepository[gorm.DB, uuid.UUID, order.Order]
}

func NewOrderRepository(baseRepo database.IBaseRepository[gorm.DB, uuid.UUID, order.Order]) order.IOrderRepository {
	return &orderRepository{
		baseRepo,
	}
}
