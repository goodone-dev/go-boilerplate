package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/BagusAK95/go-boilerplate/internal/config"
	"github.com/BagusAK95/go-boilerplate/internal/domain/customer"
	"github.com/BagusAK95/go-boilerplate/internal/domain/mail"
	"github.com/BagusAK95/go-boilerplate/internal/domain/order"
	"github.com/BagusAK95/go-boilerplate/internal/domain/product"
	"github.com/BagusAK95/go-boilerplate/internal/infrastructure/message/bus"
	httperror "github.com/BagusAK95/go-boilerplate/internal/utils/error"
	"github.com/BagusAK95/go-boilerplate/internal/utils/tracer"
	"github.com/google/uuid"
)

type OrderUsecase struct {
	customerRepo  customer.ICustomerRepository
	productRepo   product.IProductRepository
	orderRepo     order.IOrderRepository
	orderItemRepo order.IOrderItemRepository
	mailBus       bus.Bus[mail.MailSendMessage]
}

func NewOrderUsecase(
	customerRepo customer.ICustomerRepository,
	productRepo product.IProductRepository,
	orderRepo order.IOrderRepository,
	orderItemRepo order.IOrderItemRepository,
	mailBus bus.Bus[mail.MailSendMessage],
) order.IOrderUsecase {
	return &OrderUsecase{
		customerRepo:  customerRepo,
		productRepo:   productRepo,
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		mailBus:       mailBus,
	}
}

// TODO: Use DTO for return response
func (u *OrderUsecase) Create(ctx context.Context, req order.CreateOrderRequest) (res *order.Order, err error) {
	ctx, span := tracer.StartSpan(ctx, req)
	defer func() {
		span.EndSpan(err, res)
	}()

	customer, err := u.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, err
	} else if customer.ID == uuid.Nil {
		return nil, httperror.NewNotFoundError("customer not found")
	}

	var productIDs []uuid.UUID
	for _, item := range req.OrderItems {
		productIDs = append(productIDs, item.ProductID)
	}

	// TODO: Lock products
	products, err := u.productRepo.GetByIDs(ctx, productIDs)
	if err != nil {
		return nil, err
	} else if len(products) != len(req.OrderItems) {
		return nil, httperror.NewNotFoundError("one or more products not found")
	}

	productMap := make(map[uuid.UUID]product.Product)
	for _, p := range products {
		productMap[p.ID] = p
	}

	var totalAmount float64
	var orderItems []order.OrderItem

	for _, item := range req.OrderItems {
		p := productMap[item.ProductID]
		totalAmount += p.Price * float64(item.Quantity)
		orderItems = append(orderItems, order.OrderItem{
			ProductID:   p.ID,
			ProductName: p.Name,
			Quantity:    item.Quantity,
			Price:       p.Price,
			Total:       p.Price * float64(item.Quantity),
		})
	}

	tx := u.orderRepo.Begin(ctx)
	defer func() {
		if err != nil {
			u.orderRepo.Rollback(tx)
			return
		}

		u.orderRepo.Commit(tx)
	}()

	createdOrder, err := u.orderRepo.CreateWithTx(ctx, order.Order{
		CustomerID:  req.CustomerID,
		TotalAmount: totalAmount,
		Status:      "paid",
	}, tx)
	if err != nil {
		return nil, err
	}

	for i := range orderItems {
		orderItems[i].OrderID = createdOrder.ID
	}

	_, err = u.orderItemRepo.CreateBulkWithTx(ctx, orderItems, tx)
	if err != nil {
		return nil, err
	}

	u.mailBus.Publish(mail.MailSendTopic, mail.MailSendMessage{
		To:       customer.Email,
		Subject:  "Thank You for Your Purchase!",
		Template: "order_created.html",
		Data: map[string]any{
			"Name":        customer.Name,
			"OrderItems":  orderItems,
			"TotalAmount": totalAmount,
			"InvoiceURL":  fmt.Sprintf("%s/file/order/receipt/%s", config.ApplicationConfig.URL, createdOrder.ID.String()),
			"YearNow":     time.Now().Year(),
		},
	})

	return &createdOrder, nil
}
