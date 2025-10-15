package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goodone-dev/go-boilerplate/internal/config"
	"github.com/goodone-dev/go-boilerplate/internal/utils/error"
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
				c.JSON(e.Status, res)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "an unexpected error occurred"})
			}
		}
	}
}
