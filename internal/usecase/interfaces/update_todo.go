package interfaces

import (
	"context"
	"todo/internal/domain/entity"
	"todo/internal/types"

	"github.com/google/uuid"
)

type UpdateTodoUsecase interface {
	Execute(ctx context.Context, id uuid.UUID, req types.UpdateTodoRequest) (*entity.Todo, error)
}
