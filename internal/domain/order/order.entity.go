package order

import (
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
)

type Order struct {
	database.BaseEntity[uuid.UUID] `bson:",inline"`
	OrderNumber                    string    `json:"order_number" bson:"order_number"`
	CustomerID                     uuid.UUID `json:"customer_id" bson:"customer_id"`
	Total                          float64   `json:"total" bson:"total"`
	Status                         string    `json:"status" bson:"status"`
}

func (Order) TableName() string {
	return "orders"
}

func (Order) RepositoryName() string {
	return "OrderRepository"
}
