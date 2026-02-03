package order

import (
	"context"

	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderItemRepository interface {
	database.BaseRepository[gorm.DB, uuid.UUID, OrderItem]
	FindByOrderID(ctx context.Context, orderID uuid.UUID) (res []DetailOrderItem, err error)
}
