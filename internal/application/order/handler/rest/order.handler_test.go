package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/goodone-dev/go-boilerplate/internal/domain/order"
	ordermock "github.com/goodone-dev/go-boilerplate/internal/domain/order/mocks"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMain(m *testing.M) {
	logger.Disabled()
	code := m.Run()

	os.Exit(code)
}

func TestNewOrderHandler(t *testing.T) {
	mockUsecase := ordermock.NewOrderUsecaseMock(t)
	handler := NewOrderHandler(mockUsecase)

	assert.NotNil(t, handler)
}

func TestOrderHandler_Create_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := ordermock.NewOrderUsecaseMock(t)
	handler := NewOrderHandler(mockUsecase)

	customerID := uuid.New()
	productID := uuid.New()
	orderID := uuid.New()

	reqBody := order.CreateOrderRequest{
		CustomerID: customerID,
		OrderItems: []order.OrderItemRequest{
			{
				ProductID: productID,
				Quantity:  2,
			},
		},
	}

	expectedResponse := &order.CreateOrderResponse{
		ID:          orderID,
		CustomerID:  customerID,
		TotalAmount: 200.0,
		Status:      "paid",
	}

	mockUsecase.On("Create", mock.Anything, mock.AnythingOfType("order.CreateOrderRequest")).Return(expectedResponse, nil)

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.Equal(t, http.StatusOK, w.Code)
	mockUsecase.AssertExpectations(t)
}

func TestOrderHandler_Create_InvalidJSON(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := ordermock.NewOrderUsecaseMock(t)
	handler := NewOrderHandler(mockUsecase)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString("invalid json"))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.NotEmpty(t, c.Errors, "Expected errors to be set in context")
	assert.Len(t, c.Errors, 1, "Expected exactly one error")
}

func TestOrderHandler_Create_UsecaseError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := ordermock.NewOrderUsecaseMock(t)
	handler := NewOrderHandler(mockUsecase)

	customerID := uuid.New()
	productID := uuid.New()

	reqBody := order.CreateOrderRequest{
		CustomerID: customerID,
		OrderItems: []order.OrderItemRequest{
			{
				ProductID: productID,
				Quantity:  2,
			},
		},
	}

	mockUsecase.On("Create", mock.Anything, mock.AnythingOfType("order.CreateOrderRequest")).Return(nil, errors.New("usecase error"))

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.NotEmpty(t, c.Errors, "Expected errors to be set in context")
	assert.Len(t, c.Errors, 1, "Expected exactly one error")
	mockUsecase.AssertExpectations(t)
}

func TestOrderHandler_Create_ValidationError_MissingCustomerID(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := ordermock.NewOrderUsecaseMock(t)
	handler := NewOrderHandler(mockUsecase)

	productID := uuid.New()

	// Missing CustomerID (zero value UUID)
	reqBody := order.CreateOrderRequest{
		CustomerID: uuid.Nil,
		OrderItems: []order.OrderItemRequest{
			{
				ProductID: productID,
				Quantity:  2,
			},
		},
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.NotEmpty(t, c.Errors, "Expected validation errors to be set in context")
}

func TestOrderHandler_Create_ValidationError_EmptyOrderItems(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockUsecase := ordermock.NewOrderUsecaseMock(t)
	handler := NewOrderHandler(mockUsecase)

	customerID := uuid.New()

	// Empty OrderItems array
	reqBody := order.CreateOrderRequest{
		CustomerID: customerID,
		OrderItems: []order.OrderItemRequest{},
	}

	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.NotEmpty(t, c.Errors, "Expected validation errors to be set in context")
}
