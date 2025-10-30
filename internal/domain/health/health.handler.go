package health

import (
	"github.com/gin-gonic/gin"
)

type IHealthHandler interface {
	LiveCheck(c *gin.Context)
	ReadyCheck(c *gin.Context)
}
