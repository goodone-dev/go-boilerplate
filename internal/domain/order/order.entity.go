package order

import (
	database "github.com/BagusAK95/go-boilerplate/internal/infrastructure/database/sql"
	"github.com/google/uuid"
)

type Order struct {
	database.BaseEntity
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
