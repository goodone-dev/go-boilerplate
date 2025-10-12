package order

import (
	"github.com/goodonedev/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
)

type Order struct {
	database.BaseEntity[uuid.UUID] `bson:",inline"`
	CustomerID                     uuid.UUID `json:"customer_id" bson:"customer_id"`
	TotalAmount                    float64   `json:"total_amount" bson:"total_amount"`
	Status                         string    `json:"status" bson:"status"`
}

func (Order) TableName() string {
	return "orders"
}

func (Order) RepositoryName() string {
	return "OrderRepository"
}
