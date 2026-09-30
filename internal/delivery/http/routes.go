package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "todo/docs"
	"todo/internal/delivery/http/middleware"
	"todo/internal/infrastructure/container"
)

// SetupRouter registers global middleware, swagger documentation, and route groups.
func SetupRouter(log zerolog.Logger, c *container.Container) *gin.Engine {
	router := gin.New()

	router.Use(
		middleware.RequestID(),
		middleware.Logger(log),
		middleware.CORS(),
		middleware.Recovery(log),
		middleware.AuditContext(),
	)

	// Swagger documentation endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/health", c.HealthHandler.Check)

	v1 := router.Group("/api/v1")
	{
		todos := v1.Group("/todos")
		{
			todos.POST("", c.TodoHandler.Create)
			todos.GET("/:id", c.TodoHandler.GetByID)
			todos.PUT("/:id", c.TodoHandler.Update)
			todos.GET("", c.TodoHandler.List)
			todos.DELETE("/:id",c.TodoHandler.Delete)
			todos.PATCH("/:id/status", c.TodoHandler.UpdateStatus)
		}
	}

	return router
}
