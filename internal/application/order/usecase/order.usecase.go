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
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/messaging/rabbitmq/direct"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
	httperror "github.com/goodone-dev/go-boilerplate/internal/utils/http_response/error"
	"github.com/google/uuid"
)

type orderUsecase struct {
	customerRepo  customer.CustomerRepository
	productRepo   product.ProductRepository
	orderRepo     order.OrderRepository
	orderItemRepo order.OrderItemRepository
	rmqDirectPub  *direct.Publisher
}

func NewOrderUsecase(
	customerRepo customer.CustomerRepository,
	productRepo product.ProductRepository,
	orderRepo order.OrderRepository,
	orderItemRepo order.OrderItemRepository,
	rmqDirectPub *direct.Publisher,
) order.OrderUsecase {
	return &orderUsecase{
		customerRepo:  customerRepo,
		productRepo:   productRepo,
		orderRepo:     orderRepo,
		orderItemRepo: orderItemRepo,
		rmqDirectPub:  rmqDirectPub,
	}
}

func (u *orderUsecase) Create(ctx context.Context, req order.CreateOrderRequest) (res *order.CreateOrderResponse, err error) {
	ctx, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"request": req,
	})

	defer func() {
		span.SetFunctionOutput(tracer.Metadata{
			"response": res,
		}).End(err)
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

	var total float64
	var orderItems []order.OrderItem

	for _, item := range req.OrderItems {
		p := productMap[item.ProductID]
		total += p.Price * float64(item.Quantity)
		orderItems = append(orderItems, order.OrderItem{
			ProductID: p.ID,
			Quantity:  item.Quantity,
			Price:     p.Price,
			Amount:    p.Price * float64(item.Quantity),
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
		OrderNumber: fmt.Sprintf("ORD-%d", time.Now().Unix()),
		CustomerID:  req.CustomerID,
		Total:       total,
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

	err = u.rmqDirectPub.Publish(ctx, "mail.send", mail.MailSendMessage{
		To:       customer.Email,
		Subject:  "Thank You for Your Purchase!",
		Template: "order_created.html",
		Data: map[string]any{
			"Name": customer.Name,
			"OrderItems": func(oi []order.OrderItem) []map[string]any {
				var orderItems []map[string]any
				for _, item := range oi {
					orderItems = append(orderItems, map[string]any{
						"ProductName": productMap[item.ProductID].Name,
						"Quantity":    item.Quantity,
						"Price":       item.Price,
						"Amount":      item.Amount,
					})
				}
				return orderItems
			}(orderItems),
			"Total":      total,
			"InvoiceURL": fmt.Sprintf("%s/files/orders/%s/receipt", config.Application.URL, createdOrder.ID.String()),
			"YearNow":    time.Now().Year(),
		},
	})
	if err != nil {
		logger.Error(ctx, err, "❌ Failed to publish email").Write()
	}

	return &order.CreateOrderResponse{
		ID:          createdOrder.ID,
		OrderNumber: createdOrder.OrderNumber,
		CustomerID:  createdOrder.CustomerID,
		Total:       createdOrder.Total,
		Status:      createdOrder.Status,
	}, nil
}

func (u *orderUsecase) GetDetail(ctx context.Context, id uuid.UUID) (res *order.DetailOrderResponse, err error) {
	ctx, span := tracer.Start(ctx)
	span.SetFunctionInput(tracer.Metadata{
		"id": id,
	})

	defer func() {
		span.SetFunctionOutput(tracer.Metadata{
			"response": res,
		}).End(err)
	}()

	foundOrder, err := u.orderRepo.FindById(ctx, id)
	if err != nil {
		return nil, err
	} else if foundOrder == nil {
		return nil, httperror.NewNotFoundError("order with the provided ID was not found")
	}

	customer, err := u.customerRepo.FindById(ctx, foundOrder.CustomerID)
	if err != nil {
		return nil, err
	} else if customer == nil {
		return nil, httperror.NewNotFoundError("customer with the provided ID was not found")
	}

	orderItems, err := u.orderItemRepo.FindByOrderID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &order.DetailOrderResponse{
		ID:          foundOrder.ID,
		OrderNumber: foundOrder.OrderNumber,
		Customer: order.DetailOrderCustomer{
			ID:    customer.ID,
			Name:  customer.Name,
			Email: customer.Email,
		},
		OrderItems: func(orderItems []order.DetailOrderItem) []order.DetailOrderItem {
			var detailOrderItems []order.DetailOrderItem
			for _, item := range orderItems {
				detailOrderItems = append(detailOrderItems, order.DetailOrderItem{
					ProductID:   item.ProductID,
					ProductName: item.ProductName,
					Quantity:    item.Quantity,
					Price:       item.Price,
					Amount:      item.Amount,
				})
			}
			return detailOrderItems
		}(orderItems),
		Total:     foundOrder.Total,
		CreatedAt: *foundOrder.CreatedAt,
		Status:    foundOrder.Status,
	}, nil
}
