package handlers

import (
	"github.com/gin-gonic/gin"

	"todo/internal/delivery/http/response"
)

// HealthHandler reports basic liveness. It intentionally has no
// dependencies yet — once the database is wired up, it will also ping the
// DB connection here.
type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Check reports basic liveness.
// @Summary Health check
// @Description Returns the health status of the service
// @Tags health
// @Produce json
// @Success 200 {object} response.Envelope{data=map[string]any}
// @Router /health [get]
func (h *HealthHandler) Check(c *gin.Context) {
	response.Success(c, gin.H{"status": "ok"})
}
