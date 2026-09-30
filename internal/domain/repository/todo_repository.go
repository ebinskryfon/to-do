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
	Update(ctx context.Context, todo *entity.Todo) error
	List(ctx context.Context, page, pageSize int, completed *bool) ([]entity.Todo, int, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, completed bool) error
}
