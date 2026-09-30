package interfaces

import (
	"context"

	"github.com/google/uuid"
)

type DeleteTodoUsecase interface {
	Execute(ctx context.Context, id uuid.UUID) error
}