package interfaces

import (
	"context"
	"github.com/google/uuid"
	"todo/internal/domain/entity"
)

type UpdateTodoUsecase interface {
	Execute(ctx context.Context, id uuid.UUID, title, description string) (*entity.Todo, error)
}
