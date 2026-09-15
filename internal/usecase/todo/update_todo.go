package todo

import (
	"context"
	"strings"
	"time"
	"todo/internal/domain/entity"
	domainerrors "todo/internal/domain/errors"
	"todo/internal/domain/repository"
	"todo/internal/usecase/interfaces"

	"github.com/google/uuid"
)

type updateTodoUsecase struct {
	repo repository.TodoRepository
}

func NewUpdateTodoUsecase(repo repository.TodoRepository) interfaces.UpdateTodoUsecase {
	return &updateTodoUsecase{repo: repo}
}

func (uc *updateTodoUsecase) Execute(ctx context.Context, id uuid.UUID, title, description string) (*entity.Todo, error) {
	trimmedTitle := strings.TrimSpace(title)
	if trimmedTitle == "" {
		return nil, domainerrors.ErrInvalidInput
	}

	existingTodo, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	trimmedDescription := strings.TrimSpace(description)

	existingTodo.Title = trimmedTitle
	existingTodo.Description = trimmedDescription
	existingTodo.UpdatedAt = time.Now().UTC()

	if err := uc.repo.Update(ctx, existingTodo); err != nil {
		return nil, err
	}

	return existingTodo, nil

}
