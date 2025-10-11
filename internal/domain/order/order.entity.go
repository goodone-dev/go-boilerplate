package order

import (
	"github.com/BagusAK95/go-boilerplate/internal/infrastructure/database"
	"github.com/google/uuid"
)

type Order struct {
	database.BaseEntity[uuid.UUID]
	CustomerID  uuid.UUID `json:"customer_id"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
}

func (Order) TableName() string {
	return "orders"
}

func (Order) RepositoryName() string {
	return "OrderRepository"
}
