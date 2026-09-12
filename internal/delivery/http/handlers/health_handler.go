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

func (h *HealthHandler) Check(c *gin.Context) {
	response.Success(c, gin.H{"status": "ok"})
}
