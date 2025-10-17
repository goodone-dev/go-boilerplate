package order

import "github.com/google/uuid"

type CreateOrderRequest struct {
	CustomerID uuid.UUID          `json:"customer_id" validate:"required"`
	OrderItems []OrderItemRequest `json:"order_items" validate:"required,min=1,dive"`
}

type OrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" san:"min=1"`
}

type CreateOrderResponse struct {
	ID          uuid.UUID `json:"order_id"`
	CustomerID  uuid.UUID `json:"customer_id" `
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status" `
}
