package health

import (
	"github.com/gin-gonic/gin"
)

type HealthHandler interface {
	LiveCheck(c *gin.Context)
	ReadyCheck(c *gin.Context)
}
