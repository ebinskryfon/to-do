package interfaces

import (
	"context"
	"todo/internal/domain/entity"
)

// CreateTodoUsecase defines the interface for creating a new todo.
type CreateTodoUsecase interface {
	Execute(ctx context.Context, title, description string) (*entity.Todo, error)
}
