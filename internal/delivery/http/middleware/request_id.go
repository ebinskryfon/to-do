package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"
const requestIDKey = "request_id"

// RequestID assigns a unique id to every request (reusing an incoming
// X-Request-ID header if the caller already sent one), so later middleware
// and logs can correlate everything about a single request. Must run first
// in the middleware chain.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(RequestIDHeader)
		if id == "" {
			id = uuid.New().String()
		}
		c.Set(requestIDKey, id)
		c.Writer.Header().Set(RequestIDHeader, id)
		c.Next()
	}
}
