package rest

import (
	"github.com/gin-gonic/gin"
	"github.com/goodone-dev/go-boilerplate/internal/domain/order"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
	httperror "github.com/goodone-dev/go-boilerplate/internal/utils/http_response/error"
	"github.com/goodone-dev/go-boilerplate/internal/utils/http_response/success"
	"github.com/goodone-dev/go-boilerplate/internal/utils/sanitizer"
	"github.com/goodone-dev/go-boilerplate/internal/utils/validator"
)

type orderHandler struct {
	orderUsecase order.OrderUsecase
}

func NewOrderHandler(orderUsecase order.OrderUsecase) order.OrderHandler {
	return &orderHandler{
		orderUsecase: orderUsecase,
	}
}

func (h *orderHandler) Create(c *gin.Context) {
	var err error

	ctx, span := tracer.Start(c.Request.Context())
	defer func() {
		span.Stop(err)
	}()

	var req order.CreateOrderRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		c.Error(httperror.NewBadRequestError("invalid JSON payload format", err.Error()))
		return
	}

	if err = sanitizer.Sanitize(req); err != nil {
		c.Error(httperror.NewInternalServerError("failed to process request data", err.Error()))
		return
	}

	if errs := validator.Validate(req); errs != nil {
		c.Error(httperror.NewBadRequestError("request contains invalid or missing fields", errs...))
		return
	}

	order, err := h.orderUsecase.Create(ctx, req)
	if err != nil {
		c.Error(err)
		return
	}

	success.Send(c, order)
}
