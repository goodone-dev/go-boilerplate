package order

import "github.com/gin-gonic/gin"

type OrderHandler interface {
	Create(c *gin.Context)
	GetDetail(c *gin.Context)
	GetReceipt(c *gin.Context)
}
