package interfaces

import (
	"context"

	"github.com/google/uuid"
)

type UpdateTodoStatusUsecase interface {
	Execute(ctx context.Context, id uuid.UUID, completed bool) error
}