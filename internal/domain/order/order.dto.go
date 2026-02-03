package order

import "github.com/google/uuid"

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
	CustomerID  uuid.UUID `json:"customer_id"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"`
}

type DetailOrderResponse struct {
	ID          uuid.UUID           `json:"id"`
	Customer    DetailOrderCustomer `json:"customer"`
	OrderItems  []DetailOrderItem   `json:"order_items"`
	TotalAmount float64             `json:"total_amount"`
	Status      string              `json:"status"`
}

type DetailOrderItem struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int       `json:"quantity"`
	Price       float64   `json:"price"`
	SubTotal    float64   `json:"sub_total"`
}

type DetailOrderCustomer struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
}
