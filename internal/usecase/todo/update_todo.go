package todo

import (
	"context"
	"strings"
	"time"
	"todo/internal/domain/entity"
	domainerrors "todo/internal/domain/errors"
	"todo/internal/domain/repository"
	"todo/internal/types"
	"todo/internal/usecase/interfaces"

	"github.com/google/uuid"
)

type updateTodoUsecase struct {
	repo repository.TodoRepository
}

func NewUpdateTodoUsecase(repo repository.TodoRepository) interfaces.UpdateTodoUsecase {
	return &updateTodoUsecase{repo: repo}
}

func (uc *updateTodoUsecase) Execute(ctx context.Context, id uuid.UUID, req types.UpdateTodoRequest) (*entity.Todo, error) {
	trimmedTitle := strings.TrimSpace(req.Title)
	if trimmedTitle == "" {
		return nil, domainerrors.ErrInvalidInput
	}

	existingTodo, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	existingTodo.Title = trimmedTitle

	trimmedDescription := strings.TrimSpace(req.Description)
	if trimmedDescription != "" {
		existingTodo.Description = trimmedDescription
	}

	existingTodo.UpdatedAt = time.Now().UTC()

	if *req.Completed == true {
		existingTodo.Completed = true
	} else {
		existingTodo.Completed = false
	}

	if err := uc.repo.Update(ctx, existingTodo); err != nil {
		return nil, err
	}

	return existingTodo, nil

}
