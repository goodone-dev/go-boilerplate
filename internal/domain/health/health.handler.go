package health

import (
	"context"

	"github.com/gin-gonic/gin"
)

type HealthChecker interface {
	Ping(ctx context.Context) error
}

type HealthHandler interface {
	LiveCheck(c *gin.Context)
	ReadyCheck(c *gin.Context)
}
