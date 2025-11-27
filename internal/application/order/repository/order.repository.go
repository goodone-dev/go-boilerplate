package repository

import (
	"github.com/goodone-dev/go-boilerplate/internal/domain/order"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type orderRepository struct {
	database.BaseRepository[gorm.DB, uuid.UUID, order.Order]
}

func NewOrderRepository(baseRepo database.BaseRepository[gorm.DB, uuid.UUID, order.Order]) order.OrderRepository {
	return &orderRepository{
		baseRepo,
	}
}
