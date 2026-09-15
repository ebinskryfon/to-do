package repository

import (
	"context"

	"github.com/google/uuid"
	"todo/internal/domain/entity"
)

// TodoRepository defines the storage contract for Todo entities.
type TodoRepository interface {
	Create(ctx context.Context, todo *entity.Todo) error
	GetByID(ctx context.Context, id uuid.UUID) (*entity.Todo, error)
}
