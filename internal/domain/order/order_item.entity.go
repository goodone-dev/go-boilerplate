package order

import (
	"github.com/BagusAK95/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
)

type OrderItem struct {
	database.BaseEntity[uuid.UUID]
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
