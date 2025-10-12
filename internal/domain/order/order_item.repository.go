package order

import (
	"github.com/goodonedev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IOrderItemRepository interface {
	database.IBaseRepository[gorm.DB, uuid.UUID, OrderItem]
}
