package order

import (
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderItemRepository interface {
	database.BaseRepository[gorm.DB, uuid.UUID, OrderItem]
}
