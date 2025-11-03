package middleware

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/tracer"
)

type etagResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r etagResponseWriter) Write(b []byte) (int, error) {
	return r.body.Write(b)
}

// ETagMiddleware adds ETag support for GET requests
// It generates an ETag based on the response body and handles If-None-Match headers
func ETagHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		var err error

		_, span := tracer.Start(c.Request.Context())
		defer func() {
			span.Stop(err)
		}()

		// Only apply ETag for GET requests
		if c.Request.Method != http.MethodGet {
			c.Next()
			return
		}

		// Create a custom response writer to capture the response body
		blw := &etagResponseWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}

		c.Writer = blw
		c.Next()

		// Skip ETag generation if there's no body or if there was an error
		if blw.body.Len() == 0 || c.Writer.Status() >= 400 {
			_, err = blw.ResponseWriter.Write(blw.body.Bytes())
			if err != nil {
				c.Error(err)
			}
			return
		}

		// Generate ETag from response body
		hash := sha256.Sum256(blw.body.Bytes())
		etag := `"` + hex.EncodeToString(hash[:]) + `"`

		// Check If-None-Match header
		ifNoneMatch := c.GetHeader("If-None-Match")
		if ifNoneMatch == etag {
			// Content hasn't changed, return 304 Not Modified
			c.Status(http.StatusNotModified)
			return
		}

		// Set ETag header and write response
		c.Header("ETag", etag)
		_, err = blw.ResponseWriter.Write(blw.body.Bytes())
		if err != nil {
			c.Error(err)
		}
	}
}
