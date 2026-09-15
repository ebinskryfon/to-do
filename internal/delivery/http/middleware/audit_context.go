package middleware

import (
	"github.com/gin-gonic/gin"
	"todo/pkg/audit"
)

// AuditContext sets ip_address and user_agent on Gin context and Go context.Context.
func AuditContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		userAgent := c.GetHeader("User-Agent")

		c.Set("ip_address", ip)
		c.Set("user_agent", userAgent)

		auditCtx := &audit.Context{
			IP:        ip,
			UserAgent: userAgent,
		}
		ctx := audit.WithContext(c.Request.Context(), auditCtx)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
