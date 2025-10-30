package repository

import (
	"github.com/goodone-dev/go-boilerplate/internal/domain/order"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderItemRepository struct {
	database.BaseRepository[gorm.DB, uuid.UUID, order.OrderItem]
}

func NewOrderItemRepository(baseRepo database.BaseRepository[gorm.DB, uuid.UUID, order.OrderItem]) order.OrderItemRepository {
	return &orderItemRepository{
		baseRepo,
	}
}
