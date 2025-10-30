package http

import (
	"github.com/gin-gonic/gin"
	"github.com/goodone-dev/go-boilerplate/internal/domain/order"
	"github.com/goodone-dev/go-boilerplate/internal/utils/http_response/error"
	"github.com/goodone-dev/go-boilerplate/internal/utils/http_response/success"
	"github.com/goodone-dev/go-boilerplate/internal/utils/sanitizer"
	"github.com/goodone-dev/go-boilerplate/internal/utils/validator"
)

type orderHandler struct {
	orderUsecase order.IOrderUsecase
}

func NewOrderHandler(orderUsecase order.IOrderUsecase) order.IOrderHandler {
	return &orderHandler{
		orderUsecase: orderUsecase,
	}
}

func (h *orderHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var req order.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(error.NewBadRequestError("invalid JSON payload format", err.Error()))
		return
	}

	if err := sanitizer.Sanitize(req); err != nil {
		c.Error(error.NewInternalServerError("failed to process request data", err.Error()))
		return
	}

	if errs := validator.Validate(req); errs != nil {
		c.Error(error.NewBadRequestError("request contains invalid or missing fields", errs...))
		return
	}

	order, err := h.orderUsecase.Create(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	success.Send(c, order)
}
