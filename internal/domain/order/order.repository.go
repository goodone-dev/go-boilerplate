package order

import (
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type IOrderRepository interface {
	database.IBaseRepository[gorm.DB, uuid.UUID, Order]
}
