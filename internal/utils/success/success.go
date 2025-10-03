package success

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type customSuccess struct {
	Message    string `json:"message"`
	Data       any    `json:"data,omitempty"`
	Pagination any    `json:"pagination,omitempty"`
}

func Send(c *gin.Context, data any) {
	c.JSON(http.StatusOK, customSuccess{
		Message: "successfully",
		Data:    data,
	})
}
