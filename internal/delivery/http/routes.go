package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"todo/internal/delivery/http/handlers"
	"todo/internal/delivery/http/middleware"
)

// SetupRouter registers global middleware and every route group.
//
// For now this only wires the health check, since the Todo feature
// (domain/usecase/infrastructure/container) hasn't been built yet. Once it
// is, this will take a *container.Container instead of a bare
// *handlers.HealthHandler and add the /api/v1/todos routes here.
func SetupRouter(log zerolog.Logger, healthHandler *handlers.HealthHandler) *gin.Engine {
	router := gin.New()

	router.Use(
		middleware.RequestID(),
		middleware.Logger(log),
		middleware.CORS(),
		middleware.Recovery(log),
	)

	router.GET("/health", healthHandler.Check)

	return router
}
