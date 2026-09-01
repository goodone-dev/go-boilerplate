package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goodone-dev/go-boilerplate/internal/domain/order"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
	"github.com/goodone-dev/go-boilerplate/internal/utils/html"
	httperror "github.com/goodone-dev/go-boilerplate/internal/utils/http_response/error"
	"github.com/goodone-dev/go-boilerplate/internal/utils/http_response/success"
	"github.com/goodone-dev/go-boilerplate/internal/utils/sanitizer"
	"github.com/goodone-dev/go-boilerplate/internal/utils/validator"
	"github.com/google/uuid"
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
		span.End(err)
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

func (h *orderHandler) GetDetail(c *gin.Context) {
	var err error

	ctx, span := tracer.Start(c.Request.Context())
	defer func() {
		span.End(err)
	}()

	idStr := c.Param("id")
	if idStr == "" {
		c.Error(httperror.NewBadRequestError("order ID is required", ""))
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.Error(httperror.NewBadRequestError("invalid order ID format", err.Error()))
		return
	}

	order, err := h.orderUsecase.GetDetail(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	success.Send(c, order)
}

func (h *orderHandler) GetReceipt(c *gin.Context) {
	var err error

	ctx, span := tracer.Start(c.Request.Context())
	defer func() {
		span.End(err)
	}()

	idStr := c.Param("id")
	if idStr == "" {
		c.Error(httperror.NewBadRequestError("order ID is required", ""))
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		c.Error(httperror.NewBadRequestError("invalid order ID format", err.Error()))
		return
	}

	order, err := h.orderUsecase.GetDetail(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}

	c.Header("Content-Type", "text/html")
	c.Status(http.StatusOK)

	if err := html.ExecuteTemplate(c.Writer, "order_receipt.html", order); err != nil {
		c.Error(httperror.NewInternalServerError("failed to generate receipt", err.Error()))
		return
	}
}
