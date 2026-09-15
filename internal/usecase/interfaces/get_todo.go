package interfaces

import (
	"context"

	"github.com/google/uuid"
	"todo/internal/domain/entity"
)

type GetTodoUsecase interface {
	Execute(ctx context.Context, id uuid.UUID) (*entity.Todo, error)
}
