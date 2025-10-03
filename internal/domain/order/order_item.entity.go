package order

import (
	database "github.com/BagusAK95/go-skeleton/internal/infrastructure/database/sql"
	"github.com/google/uuid"
)

type OrderItem struct {
	database.BaseEntity
	OrderID     uuid.UUID `json:"order_id"`
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name" gorm:"-"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	Total       float64   `json:"total" gorm:"-"`
}

func (OrderItem) TableName() string {
	return "order_items"
}

func (OrderItem) RepositoryName() string {
	return "OrderItemRepository"
}
