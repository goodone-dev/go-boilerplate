package middleware

import (
	"net/http"

	"github.com/BagusAK95/go-skeleton/internal/utils/error"
	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if e, ok := err.(*error.CustomError); ok {
				res := gin.H{"message": e.Message}
				if len(e.Errors) > 0 {
					res["errors"] = e.Errors
				}
				c.JSON(e.Status, res)
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "an unexpected error occurred"})
			}
		}
	}
}
