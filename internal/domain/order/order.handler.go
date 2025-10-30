package order

import "github.com/gin-gonic/gin"

type IOrderHandler interface {
	Create(c *gin.Context)
}
