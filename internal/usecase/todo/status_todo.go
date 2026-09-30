package todo

import (
	"context"
	"todo/internal/domain/repository"
	"todo/internal/usecase/interfaces"

	"github.com/google/uuid"
)

type updateTodoStatusUsecase struct {
	repo repository.TodoRepository
}

func NewUpdateTodoStatusUsecase(repo repository.TodoRepository) interfaces.UpdateTodoStatusUsecase {
	return &updateTodoStatusUsecase{repo: repo}
}

func (uc *updateTodoStatusUsecase) Execute(ctx context.Context, id uuid.UUID, completed bool) error {
	err := uc.repo.UpdateStatus(ctx, id, completed)
	if err != nil {
		return err
	}

	return nil
}