package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"todo/internal/delivery/http/response"
)

// Recovery recovers from a panic in any handler further down the chain,
// logs it, and returns a clean 500 envelope instead of crashing the process
// or leaking Gin's default plaintext panic output. Must run last in the
// middleware chain so it wraps every handler.
func Recovery(log zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Error().
					Str("request_id", c.GetString(requestIDKey)).
					Interface("panic", err).
					Msg("recovered from panic")
				response.Error(c, http.StatusInternalServerError, "INTERNAL_ERROR", "unexpected error")
			}
		}()
		c.Next()
	}
}
