package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/domain/customer"
	"github.com/goodone-dev/go-boilerplate/internal/domain/mail"
	"github.com/goodone-dev/go-boilerplate/internal/domain/order"
	"github.com/goodone-dev/go-boilerplate/internal/domain/product"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/message/bus"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
	httperror "github.com/goodone-dev/go-boilerplate/internal/utils/http_response/error"
	"github.com/google/uuid"
)

type orderUsecase struct {
	customerRepo  customer.CustomerRepository
	productRepo   product.ProductRepository
	orderRepo     order.OrderRepository
	orderItemRepo order.OrderItemRepository
	mailBus       bus.Bus[mail.MailSendMessage]
}

func NewOrderUsecase(
	customerRepo customer.CustomerRepository,
	productRepo product.ProductRepository,
	orderRepo order.OrderRepository,
	orderItemRepo order.OrderItemRepository,
	mailBus bus.Bus[mail.MailSendMessage],
) order.OrderUsecase {
	return &orderUsecase{
		customerRepo:  customerRepo,
		productRepo:   productRepo,
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		mailBus:       mailBus,
	}
}

func (u *orderUsecase) Create(ctx context.Context, req order.CreateOrderRequest) (res *order.CreateOrderResponse, err error) {
	ctx, span := tracer.Start(ctx, req)
	defer func() {
		span.Stop(err, res)
	}()

	customer, err := u.customerRepo.FindById(ctx, req.CustomerID)
	if err != nil {
		return nil, err
	} else if customer == nil {
		return nil, httperror.NewNotFoundError("customer with the provided ID was not found")
	}

	var productIDs []uuid.UUID
	for _, item := range req.OrderItems {
		productIDs = append(productIDs, item.ProductID)
	}

	// TODO: Lock products
	products, err := u.productRepo.FindByIds(ctx, productIDs)
	if err != nil {
		return nil, err
	} else if len(products) != len(req.OrderItems) {
		return nil, httperror.NewNotFoundError("one or more requested products could not be found")
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

	trx, err := u.orderRepo.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			u.orderRepo.Rollback(trx)
			return
		}

		u.orderRepo.Commit(trx)
	}()

	createdOrder, err := u.orderRepo.Insert(ctx, order.Order{
		CustomerID:  req.CustomerID,
		TotalAmount: totalAmount,
		Status:      "paid",
	}, trx)
	if err != nil {
		return nil, err
	}

	for i := range orderItems {
		orderItems[i].OrderID = createdOrder.ID
	}

	_, err = u.orderItemRepo.InsertMany(ctx, orderItems, trx)
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

	return &order.CreateOrderResponse{
		ID:          createdOrder.ID,
		CustomerID:  createdOrder.CustomerID,
		TotalAmount: createdOrder.TotalAmount,
		Status:      createdOrder.Status,
	}, nil
}
