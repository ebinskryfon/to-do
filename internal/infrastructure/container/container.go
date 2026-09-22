package container

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"todo/internal/delivery/http/handlers"
	"todo/internal/domain/repository"
	"todo/internal/infrastructure/persistence"
	"todo/internal/usecase/interfaces"
	"todo/internal/usecase/todo"
)

// Container manages dependency injection for the entire application.
type Container struct {
	// Repositories
	TodoRepo repository.TodoRepository

	// Usecases
	CreateTodoUC      interfaces.CreateTodoUsecase
	GetTodoUsecase    interfaces.GetTodoUsecase
	UpdateTodoUsecase interfaces.UpdateTodoUsecase
	ListTodoUsecase	  interfaces.ListTodoUsecase

	// Handlers
	TodoHandler   *handlers.TodoHandler
	HealthHandler *handlers.HealthHandler
}

// NewContainer initializes and wires dependencies bottom-up (DB -> Repo -> Usecase -> Handler).
func NewContainer(db *gorm.DB, log zerolog.Logger) *Container {
	// 1. Persistence Layer
	todoRepo := persistence.NewTodoRepository(db)

	// 2. Usecase Layer
	createTodoUC := todo.NewCreateTodoUsecase(todoRepo)
	getTodoUC := todo.NewGetTodoUsecase(todoRepo)
	updateTodoUC := todo.NewUpdateTodoUsecase(todoRepo)
	listTodoUC := todo.NewListTodoUsecase(todoRepo)

	// 3. Delivery / Handler Layer
	todoHandler := handlers.NewTodoHandler(createTodoUC, getTodoUC, updateTodoUC, listTodoUC, log)
	healthHandler := handlers.NewHealthHandler()

	return &Container{
		TodoRepo:          todoRepo,
		CreateTodoUC:      createTodoUC,
		GetTodoUsecase:    getTodoUC,
		UpdateTodoUsecase: updateTodoUC,
		ListTodoUsecase: listTodoUC,
		TodoHandler:       todoHandler,
		HealthHandler:     healthHandler,
	}
}
