package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
	"github.com/goodone-dev/go-boilerplate/internal/utils/http_response/error"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if e, ok := err.(*error.CustomError); ok {
				res := gin.H{"message": e.Message}
				if len(e.Errors) > 0 && config.ApplicationConfig.Env != config.EnvProd {
					res["errors"] = e.Errors
				}

				logger.Errorf(c.Request.Context(), errors.New(strings.Join(e.Errors, ", ")), "❌ Failed to process %s %s: %s", c.Request.Method, c.Request.URL.Path, e.Message)

				c.JSON(e.Status, res)
			} else {
				res := gin.H{"message": "an unexpected error occurred"}
				if config.ApplicationConfig.Env != config.EnvProd {
					res["errors"] = []string{err.Error()}
				}

				logger.Errorf(c.Request.Context(), err, "❌ Failed to process %s %s: An unexpected error occurred", c.Request.Method, c.Request.URL.Path)

				c.JSON(http.StatusInternalServerError, res)
			}
		}
	}
}
