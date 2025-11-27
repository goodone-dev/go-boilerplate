package http

import (
	"net/http"
	"reflect"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/goodone-dev/go-boilerplate/internal/domain/health"
	httperror "github.com/goodone-dev/go-boilerplate/internal/utils/http_response/error"
)

type healthHandler struct {
	checkers []health.HealthChecker
}

func NewHealthHandler(checkers ...health.HealthChecker) health.HealthHandler {
	return &healthHandler{
		checkers: checkers,
	}
}

func (h *healthHandler) LiveCheck(c *gin.Context) {
	c.JSON(http.StatusOK, health.HealthResponse{Status: health.StatusUp})
}

func (h *healthHandler) ReadyCheck(c *gin.Context) {
	ctx := c.Request.Context()

	res := make(map[string]health.HealthResponse)
	for _, checker := range h.checkers {
		packageName := parsePackageName(checker)

		if err := checker.Ping(ctx); err != nil {
			c.Error(httperror.NewServiceUnavailableError("service dependency health check failed", err.Error()))
			return
		}

		res[packageName] = health.HealthResponse{Status: health.StatusUp}
	}

	c.JSON(http.StatusOK, res)
}

func parsePackageName(checker health.HealthChecker) string {
	n := reflect.TypeOf(checker).String()
	r := regexp.MustCompile(`\*?([^.]+)`)

	matches := r.FindStringSubmatch(n)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}
