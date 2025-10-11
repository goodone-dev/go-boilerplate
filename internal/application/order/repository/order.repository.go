package repository

import (
	"github.com/BagusAK95/go-boilerplate/internal/domain/order"
	"github.com/BagusAK95/go-boilerplate/internal/infrastructure/database"
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
