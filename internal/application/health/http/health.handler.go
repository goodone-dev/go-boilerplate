package http

import (
	"net/http"
	"reflect"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/goodone-dev/go-boilerplate/internal/domain/health"
)

type HealthHandler struct {
	services []health.IHealthService
}

func NewHealthHandler(services ...health.IHealthService) *HealthHandler {
	return &HealthHandler{
		services: services,
	}
}

func (h *HealthHandler) HealthLiveCheck(c *gin.Context) {
	c.JSON(http.StatusOK, health.HealthStatus{Status: health.StatusUp})
}

func (h *HealthHandler) HealthReadyCheck(c *gin.Context) {
	ctx := c.Request.Context()

	res := make(map[string]health.HealthStatus)
	for _, service := range h.services {
		packageName := parsePackageName(service)

		if err := service.Ping(ctx); err != nil {
			res[packageName] = health.HealthStatus{Status: health.StatusDown}

			c.JSON(http.StatusServiceUnavailable, res)
			return
		}

		res[packageName] = health.HealthStatus{Status: health.StatusUp}
	}

	c.JSON(http.StatusOK, res)
}

func parsePackageName(service health.IHealthService) string {
	n := reflect.TypeOf(service).String()
	r := regexp.MustCompile(`\*?([^.]+)`)

	matches := r.FindStringSubmatch(n)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}
