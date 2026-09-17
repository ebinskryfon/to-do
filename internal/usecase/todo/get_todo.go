package todo

import (
	"context"

	"github.com/google/uuid"
	"todo/internal/domain/entity"
	"todo/internal/domain/repository"
	"todo/internal/usecase/interfaces"
)

type getTodoUsecase struct {
	repo repository.TodoRepository
}

// NewGetTodoUsecase creates a new GetTodoUsecase implementation.
func NewGetTodoUsecase(repo repository.TodoRepository) interfaces.GetTodoUsecase {
	return &getTodoUsecase{repo: repo}
}

func (uc *getTodoUsecase) Execute(ctx context.Context, id uuid.UUID) (*entity.Todo, error) {
	return uc.repo.GetByID(ctx, id)
}
