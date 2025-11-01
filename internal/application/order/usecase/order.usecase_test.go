package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/goodone-dev/go-boilerplate/internal/domain/customer"
	customermock "github.com/goodone-dev/go-boilerplate/internal/domain/customer/mocks"
	"github.com/goodone-dev/go-boilerplate/internal/domain/mail"
	"github.com/goodone-dev/go-boilerplate/internal/domain/order"
	ordermock "github.com/goodone-dev/go-boilerplate/internal/domain/order/mocks"
	"github.com/goodone-dev/go-boilerplate/internal/domain/product"
	productmock "github.com/goodone-dev/go-boilerplate/internal/domain/product/mocks"
	busmock "github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/bus/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestNewOrderUsecase(t *testing.T) {
	mockCustomerRepo := customermock.NewCustomerRepositoryMock(t)
	mockProductRepo := productmock.NewProductRepositoryMock(t)
	mockOrderRepo := ordermock.NewOrderRepositoryMock(t)
	mockOrderItemRepo := ordermock.NewOrderItemRepositoryMock(t)
	mockMailBus := busmock.NewBusMock[mail.MailSendMessage](t)

	usecase := NewOrderUsecase(
		mockCustomerRepo,
		mockProductRepo,
		mockOrderRepo,
		mockOrderItemRepo,
		mockMailBus,
	)

	assert.NotNil(t, usecase)
}

func TestOrderUsecase_Create_Success(t *testing.T) {
	// Setup
	ctx := context.Background()
	customerID := uuid.New()
	productID1 := uuid.New()
	productID2 := uuid.New()
	orderID := uuid.New()

	mockCustomerRepo := customermock.NewCustomerRepositoryMock(t)
	mockProductRepo := productmock.NewProductRepositoryMock(t)
	mockOrderRepo := ordermock.NewOrderRepositoryMock(t)
	mockOrderItemRepo := ordermock.NewOrderItemRepositoryMock(t)
	mockMailBus := busmock.NewBusMock[mail.MailSendMessage](t)

	// Mock data
	mockCustomer := &customer.Customer{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockCustomer.ID = customerID

	mockProducts := []product.Product{
		{
			Name:  "Product 1",
			Price: 100.0,
		},
		{
			Name:  "Product 2",
			Price: 200.0,
		},
	}
	mockProducts[0].ID = productID1
	mockProducts[1].ID = productID2

	mockOrder := order.Order{
		CustomerID:  customerID,
		TotalAmount: 500.0,
		Status:      "paid",
	}
	mockOrder.ID = orderID

	req := order.CreateOrderRequest{
		CustomerID: customerID,
		OrderItems: []order.OrderItemRequest{
			{
				ProductID: productID1,
				Quantity:  2,
			},
			{
				ProductID: productID2,
				Quantity:  1,
			},
		},
	}

	mockTrx := &gorm.DB{}

	// Mock expectations
	mockCustomerRepo.EXPECT().FindById(ctx, customerID).Return(mockCustomer, nil)
	mockProductRepo.EXPECT().FindByIds(ctx, []uuid.UUID{productID1, productID2}).Return(mockProducts, nil)
	mockOrderRepo.EXPECT().Begin(ctx).Return(mockTrx, nil)
	mockOrderRepo.EXPECT().Insert(ctx, mock.MatchedBy(func(o order.Order) bool {
		return o.CustomerID == customerID && o.TotalAmount == 400.0 && o.Status == "paid"
	}), mockTrx).Return(mockOrder, nil)
	mockOrderItemRepo.EXPECT().InsertMany(ctx, mock.MatchedBy(func(items []order.OrderItem) bool {
		return len(items) == 2 && items[0].OrderID == orderID && items[1].OrderID == orderID
	}), mockTrx).Return([]order.OrderItem{}, nil)
	mockOrderRepo.EXPECT().Commit(mockTrx).Return(mockTrx)
	mockMailBus.EXPECT().Publish(mail.MailSendTopic, mock.MatchedBy(func(msg mail.MailSendMessage) bool {
		return msg.To == "john@example.com" && msg.Subject == "Thank You for Your Purchase!"
	})).Return()

	// Execute
	usecase := NewOrderUsecase(
		mockCustomerRepo,
		mockProductRepo,
		mockOrderRepo,
		mockOrderItemRepo,
		mockMailBus,
	)

	result, err := usecase.Create(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, orderID, result.ID)
	assert.Equal(t, customerID, result.CustomerID)
	assert.Equal(t, 500.0, result.TotalAmount)
	assert.Equal(t, "paid", result.Status)
}

func TestOrderUsecase_Create_CustomerNotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()

	mockCustomerRepo := customermock.NewCustomerRepositoryMock(t)
	mockProductRepo := productmock.NewProductRepositoryMock(t)
	mockOrderRepo := ordermock.NewOrderRepositoryMock(t)
	mockOrderItemRepo := ordermock.NewOrderItemRepositoryMock(t)
	mockMailBus := busmock.NewBusMock[mail.MailSendMessage](t)

	req := order.CreateOrderRequest{
		CustomerID: customerID,
		OrderItems: []order.OrderItemRequest{
			{
				ProductID: productID,
				Quantity:  1,
			},
		},
	}

	// Mock expectations - customer not found
	mockCustomerRepo.EXPECT().FindById(ctx, customerID).Return(nil, nil)

	// Execute
	usecase := NewOrderUsecase(
		mockCustomerRepo,
		mockProductRepo,
		mockOrderRepo,
		mockOrderItemRepo,
		mockMailBus,
	)

	result, err := usecase.Create(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "customer with the provided ID was not found")
}

func TestOrderUsecase_Create_CustomerRepoError(t *testing.T) {
	// Setup
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()

	mockCustomerRepo := customermock.NewCustomerRepositoryMock(t)
	mockProductRepo := productmock.NewProductRepositoryMock(t)
	mockOrderRepo := ordermock.NewOrderRepositoryMock(t)
	mockOrderItemRepo := ordermock.NewOrderItemRepositoryMock(t)
	mockMailBus := busmock.NewBusMock[mail.MailSendMessage](t)

	req := order.CreateOrderRequest{
		CustomerID: customerID,
		OrderItems: []order.OrderItemRequest{
			{
				ProductID: productID,
				Quantity:  1,
			},
		},
	}

	expectedError := errors.New("database connection error")

	// Mock expectations
	mockCustomerRepo.EXPECT().FindById(ctx, customerID).Return(nil, expectedError)

	// Execute
	usecase := NewOrderUsecase(
		mockCustomerRepo,
		mockProductRepo,
		mockOrderRepo,
		mockOrderItemRepo,
		mockMailBus,
	)

	result, err := usecase.Create(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

func TestOrderUsecase_Create_ProductsNotFound(t *testing.T) {
	// Setup
	ctx := context.Background()
	customerID := uuid.New()
	productID1 := uuid.New()
	productID2 := uuid.New()

	mockCustomerRepo := customermock.NewCustomerRepositoryMock(t)
	mockProductRepo := productmock.NewProductRepositoryMock(t)
	mockOrderRepo := ordermock.NewOrderRepositoryMock(t)
	mockOrderItemRepo := ordermock.NewOrderItemRepositoryMock(t)
	mockMailBus := busmock.NewBusMock[mail.MailSendMessage](t)

	mockCustomer := &customer.Customer{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockCustomer.ID = customerID

	req := order.CreateOrderRequest{
		CustomerID: customerID,
		OrderItems: []order.OrderItemRequest{
			{
				ProductID: productID1,
				Quantity:  1,
			},
			{
				ProductID: productID2,
				Quantity:  1,
			},
		},
	}

	// Return only one product when two are requested
	mockProducts := []product.Product{
		{
			Name:  "Product 1",
			Price: 100.0,
		},
	}
	mockProducts[0].ID = productID1

	// Mock expectations
	mockCustomerRepo.EXPECT().FindById(ctx, customerID).Return(mockCustomer, nil)
	mockProductRepo.EXPECT().FindByIds(ctx, []uuid.UUID{productID1, productID2}).Return(mockProducts, nil)

	// Execute
	usecase := NewOrderUsecase(
		mockCustomerRepo,
		mockProductRepo,
		mockOrderRepo,
		mockOrderItemRepo,
		mockMailBus,
	)

	result, err := usecase.Create(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "one or more requested products could not be found")
}

func TestOrderUsecase_Create_ProductRepoError(t *testing.T) {
	// Setup
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()

	mockCustomerRepo := customermock.NewCustomerRepositoryMock(t)
	mockProductRepo := productmock.NewProductRepositoryMock(t)
	mockOrderRepo := ordermock.NewOrderRepositoryMock(t)
	mockOrderItemRepo := ordermock.NewOrderItemRepositoryMock(t)
	mockMailBus := busmock.NewBusMock[mail.MailSendMessage](t)

	mockCustomer := &customer.Customer{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockCustomer.ID = customerID

	req := order.CreateOrderRequest{
		CustomerID: customerID,
		OrderItems: []order.OrderItemRequest{
			{
				ProductID: productID,
				Quantity:  1,
			},
		},
	}

	expectedError := errors.New("product database error")

	// Mock expectations
	mockCustomerRepo.EXPECT().FindById(ctx, customerID).Return(mockCustomer, nil)
	mockProductRepo.EXPECT().FindByIds(ctx, []uuid.UUID{productID}).Return(nil, expectedError)

	// Execute
	usecase := NewOrderUsecase(
		mockCustomerRepo,
		mockProductRepo,
		mockOrderRepo,
		mockOrderItemRepo,
		mockMailBus,
	)

	result, err := usecase.Create(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

func TestOrderUsecase_Create_BeginTransactionError(t *testing.T) {
	// Setup
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()

	mockCustomerRepo := customermock.NewCustomerRepositoryMock(t)
	mockProductRepo := productmock.NewProductRepositoryMock(t)
	mockOrderRepo := ordermock.NewOrderRepositoryMock(t)
	mockOrderItemRepo := ordermock.NewOrderItemRepositoryMock(t)
	mockMailBus := busmock.NewBusMock[mail.MailSendMessage](t)

	mockCustomer := &customer.Customer{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockCustomer.ID = customerID

	mockProducts := []product.Product{
		{
			Name:  "Product 1",
			Price: 100.0,
		},
	}
	mockProducts[0].ID = productID

	req := order.CreateOrderRequest{
		CustomerID: customerID,
		OrderItems: []order.OrderItemRequest{
			{
				ProductID: productID,
				Quantity:  1,
			},
		},
	}

	expectedError := errors.New("transaction begin error")

	// Mock expectations
	mockCustomerRepo.EXPECT().FindById(ctx, customerID).Return(mockCustomer, nil)
	mockProductRepo.EXPECT().FindByIds(ctx, []uuid.UUID{productID}).Return(mockProducts, nil)
	mockOrderRepo.EXPECT().Begin(ctx).Return(nil, expectedError)

	// Execute
	usecase := NewOrderUsecase(
		mockCustomerRepo,
		mockProductRepo,
		mockOrderRepo,
		mockOrderItemRepo,
		mockMailBus,
	)

	result, err := usecase.Create(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

func TestOrderUsecase_Create_InsertOrderError(t *testing.T) {
	// Setup
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()

	mockCustomerRepo := customermock.NewCustomerRepositoryMock(t)
	mockProductRepo := productmock.NewProductRepositoryMock(t)
	mockOrderRepo := ordermock.NewOrderRepositoryMock(t)
	mockOrderItemRepo := ordermock.NewOrderItemRepositoryMock(t)
	mockMailBus := busmock.NewBusMock[mail.MailSendMessage](t)

	mockCustomer := &customer.Customer{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockCustomer.ID = customerID

	mockProducts := []product.Product{
		{
			Name:  "Product 1",
			Price: 100.0,
		},
	}
	mockProducts[0].ID = productID

	req := order.CreateOrderRequest{
		CustomerID: customerID,
		OrderItems: []order.OrderItemRequest{
			{
				ProductID: productID,
				Quantity:  2,
			},
		},
	}

	mockTrx := &gorm.DB{}
	expectedError := errors.New("insert order error")

	// Mock expectations
	mockCustomerRepo.EXPECT().FindById(ctx, customerID).Return(mockCustomer, nil)
	mockProductRepo.EXPECT().FindByIds(ctx, []uuid.UUID{productID}).Return(mockProducts, nil)
	mockOrderRepo.EXPECT().Begin(ctx).Return(mockTrx, nil)
	mockOrderRepo.EXPECT().Insert(ctx, mock.Anything, mockTrx).Return(order.Order{}, expectedError)
	mockOrderRepo.EXPECT().Rollback(mockTrx).Return(mockTrx)

	// Execute
	usecase := NewOrderUsecase(
		mockCustomerRepo,
		mockProductRepo,
		mockOrderRepo,
		mockOrderItemRepo,
		mockMailBus,
	)

	result, err := usecase.Create(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

func TestOrderUsecase_Create_InsertOrderItemsError(t *testing.T) {
	// Setup
	ctx := context.Background()
	customerID := uuid.New()
	productID := uuid.New()
	orderID := uuid.New()

	mockCustomerRepo := customermock.NewCustomerRepositoryMock(t)
	mockProductRepo := productmock.NewProductRepositoryMock(t)
	mockOrderRepo := ordermock.NewOrderRepositoryMock(t)
	mockOrderItemRepo := ordermock.NewOrderItemRepositoryMock(t)
	mockMailBus := busmock.NewBusMock[mail.MailSendMessage](t)

	mockCustomer := &customer.Customer{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockCustomer.ID = customerID

	mockProducts := []product.Product{
		{
			Name:  "Product 1",
			Price: 100.0,
		},
	}
	mockProducts[0].ID = productID

	mockOrder := order.Order{
		CustomerID:  customerID,
		TotalAmount: 200.0,
		Status:      "paid",
	}
	mockOrder.ID = orderID

	req := order.CreateOrderRequest{
		CustomerID: customerID,
		OrderItems: []order.OrderItemRequest{
			{
				ProductID: productID,
				Quantity:  2,
			},
		},
	}

	mockTrx := &gorm.DB{}
	expectedError := errors.New("insert order items error")

	// Mock expectations
	mockCustomerRepo.EXPECT().FindById(ctx, customerID).Return(mockCustomer, nil)
	mockProductRepo.EXPECT().FindByIds(ctx, []uuid.UUID{productID}).Return(mockProducts, nil)
	mockOrderRepo.EXPECT().Begin(ctx).Return(mockTrx, nil)
	mockOrderRepo.EXPECT().Insert(ctx, mock.Anything, mockTrx).Return(mockOrder, nil)
	mockOrderItemRepo.EXPECT().InsertMany(ctx, mock.Anything, mockTrx).Return(nil, expectedError)
	mockOrderRepo.EXPECT().Rollback(mockTrx).Return(mockTrx)

	// Execute
	usecase := NewOrderUsecase(
		mockCustomerRepo,
		mockProductRepo,
		mockOrderRepo,
		mockOrderItemRepo,
		mockMailBus,
	)

	result, err := usecase.Create(ctx, req)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedError, err)
}

func TestOrderUsecase_Create_CalculatesTotalAmountCorrectly(t *testing.T) {
	// Setup
	ctx := context.Background()
	customerID := uuid.New()
	productID1 := uuid.New()
	productID2 := uuid.New()
	productID3 := uuid.New()
	orderID := uuid.New()

	mockCustomerRepo := customermock.NewCustomerRepositoryMock(t)
	mockProductRepo := productmock.NewProductRepositoryMock(t)
	mockOrderRepo := ordermock.NewOrderRepositoryMock(t)
	mockOrderItemRepo := ordermock.NewOrderItemRepositoryMock(t)
	mockMailBus := busmock.NewBusMock[mail.MailSendMessage](t)

	mockCustomer := &customer.Customer{
		Name:  "John Doe",
		Email: "john@example.com",
	}
	mockCustomer.ID = customerID

	mockProducts := []product.Product{
		{
			Name:  "Product 1",
			Price: 50.0,
		},
		{
			Name:  "Product 2",
			Price: 75.5,
		},
		{
			Name:  "Product 3",
			Price: 120.25,
		},
	}
	mockProducts[0].ID = productID1
	mockProducts[1].ID = productID2
	mockProducts[2].ID = productID3

	mockOrder := order.Order{
		CustomerID:  customerID,
		TotalAmount: 421.25, // (50*3) + (75.5*2) + (120.25*1) = 150 + 151 + 120.25
		Status:      "paid",
	}
	mockOrder.ID = orderID

	req := order.CreateOrderRequest{
		CustomerID: customerID,
		OrderItems: []order.OrderItemRequest{
			{
				ProductID: productID1,
				Quantity:  3,
			},
			{
				ProductID: productID2,
				Quantity:  2,
			},
			{
				ProductID: productID3,
				Quantity:  1,
			},
		},
	}

	mockTrx := &gorm.DB{}

	// Mock expectations
	mockCustomerRepo.EXPECT().FindById(ctx, customerID).Return(mockCustomer, nil)
	mockProductRepo.EXPECT().FindByIds(ctx, []uuid.UUID{productID1, productID2, productID3}).Return(mockProducts, nil)
	mockOrderRepo.EXPECT().Begin(ctx).Return(mockTrx, nil)
	mockOrderRepo.EXPECT().Insert(ctx, mock.MatchedBy(func(o order.Order) bool {
		// Verify total amount calculation: (50*3) + (75.5*2) + (120.25*1) = 421.25
		return o.TotalAmount == 421.25
	}), mockTrx).Return(mockOrder, nil)
	mockOrderItemRepo.EXPECT().InsertMany(ctx, mock.MatchedBy(func(items []order.OrderItem) bool {
		if len(items) != 3 {
			return false
		}
		// Verify individual item totals
		return items[0].Total == 150.0 && items[1].Total == 151.0 && items[2].Total == 120.25
	}), mockTrx).Return([]order.OrderItem{}, nil)
	mockOrderRepo.EXPECT().Commit(mockTrx).Return(mockTrx)
	mockMailBus.EXPECT().Publish(mail.MailSendTopic, mock.Anything).Return()

	// Execute
	usecase := NewOrderUsecase(
		mockCustomerRepo,
		mockProductRepo,
		mockOrderRepo,
		mockOrderItemRepo,
		mockMailBus,
	)

	result, err := usecase.Create(ctx, req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 421.25, result.TotalAmount)
}
