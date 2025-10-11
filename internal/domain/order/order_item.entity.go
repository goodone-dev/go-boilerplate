package order

import (
	"github.com/BagusAK95/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
)

type OrderItem struct {
	database.BaseEntity[uuid.UUID] `bson:",inline"`
	OrderID                        uuid.UUID `json:"order_id" bson:"order_id"`
	ProductID                      uuid.UUID `json:"product_id" bson:"product_id"`
	ProductName                    string    `json:"product_name" gorm:"-" bson:"product_name"`
	Quantity                       int       `json:"quantity" bson:"quantity"`
	Price                          float64   `json:"price" bson:"price"`
	Total                          float64   `json:"total" gorm:"-" bson:"total"`
}

func (OrderItem) TableName() string {
	return "order_items"
}

func (OrderItem) RepositoryName() string {
	return "OrderItemRepository"
}
