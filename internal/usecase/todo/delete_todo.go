package todo

import (
	"context"
	"todo/internal/domain/repository"
	"todo/internal/usecase/interfaces"

	"github.com/google/uuid"
)

type deleteTodoUsecase struct {
	repo repository.TodoRepository
}

func NewDeleteTodoUsecase(repo repository.TodoRepository) interfaces.DeleteTodoUsecase {
	return &deleteTodoUsecase{repo: repo}
}

func (uc *deleteTodoUsecase) Execute(ctx context.Context, id uuid.UUID) error {
	err := uc.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	return nil
}