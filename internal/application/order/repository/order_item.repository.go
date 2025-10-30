package repository

import (
	"github.com/goodone-dev/go-boilerplate/internal/domain/order"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderItemRepository struct {
	database.IBaseRepository[gorm.DB, uuid.UUID, order.OrderItem]
}

func NewOrderItemRepository(baseRepo database.IBaseRepository[gorm.DB, uuid.UUID, order.OrderItem]) order.IOrderItemRepository {
	return &orderItemRepository{
		baseRepo,
	}
}
