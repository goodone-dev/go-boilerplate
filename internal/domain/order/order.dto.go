package order

import (
	"time"

	"github.com/google/uuid"
)

type CreateOrderRequest struct {
	CustomerID uuid.UUID          `json:"customer_id" validate:"required"`
	OrderItems []OrderItemRequest `json:"order_items" validate:"required,min=1,dive"`
}

type OrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

type CreateOrderResponse struct {
	ID          uuid.UUID `json:"id"`
	OrderNumber string    `json:"order_number"`
	CustomerID  uuid.UUID `json:"customer_id"`
	Total       float64   `json:"total"`
	Status      string    `json:"status"`
}

type DetailOrderResponse struct {
	ID          uuid.UUID           `json:"id"`
	OrderNumber string              `json:"order_number"`
	Customer    DetailOrderCustomer `json:"customer"`
	OrderItems  []DetailOrderItem   `json:"order_items"`
	Total       float64             `json:"total"`
	CreatedAt   time.Time           `json:"created_at"`
	Status      string              `json:"status"`
}

type DetailOrderItem struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	Amount      float64   `json:"amount"`
}

type DetailOrderCustomer struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}
