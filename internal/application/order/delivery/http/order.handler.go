package http

import (
	"github.com/BagusAK95/go-skeleton/internal/domain/order"
	"github.com/BagusAK95/go-skeleton/internal/utils/error"
	"github.com/BagusAK95/go-skeleton/internal/utils/success"
	"github.com/BagusAK95/go-skeleton/internal/utils/validator"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderUsecase order.IOrderUsecase
}

func NewOrderHandler(orderUsecase order.IOrderUsecase) *OrderHandler {
	return &OrderHandler{
		orderUsecase: orderUsecase,
	}
}

func (h *OrderHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var req order.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(error.NewInternalServerError(err.Error()))
		return
	}

	if errs := validator.Validate(req); errs != nil {
		c.Error(error.NewBadRequestError("invalid request body", errs...))
		return
	}

	order, err := h.orderUsecase.Create(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	success.Send(c, order)
}
