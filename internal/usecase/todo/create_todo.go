package todo

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"todo/internal/domain/entity"
	domainerrors "todo/internal/domain/errors"
	"todo/internal/domain/repository"
	"todo/internal/usecase/interfaces"
)

type createTodoUsecase struct {
	repo repository.TodoRepository
}

// NewCreateTodoUsecase creates a new CreateTodoUsecase implementation.
func NewCreateTodoUsecase(repo repository.TodoRepository) interfaces.CreateTodoUsecase {
	return &createTodoUsecase{repo: repo}
}

func (uc *createTodoUsecase) Execute(ctx context.Context, title, description string) (*entity.Todo, error) {
	trimmedTitle := strings.TrimSpace(title)
	if trimmedTitle == "" {
		return nil, domainerrors.ErrInvalidInput
	}

	now := time.Now().UTC()
	todo := &entity.Todo{
		SoftDeletableEntity: entity.SoftDeletableEntity{
			AuditableEntity: entity.AuditableEntity{
				BaseModel: entity.BaseModel{
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			IsActive: true,
		},
		ID:          uuid.New(),
		Title:       trimmedTitle,
		Description: strings.TrimSpace(description),
		Completed:   false,
	}

	if err := uc.repo.Create(ctx, todo); err != nil {
		return nil, err
	}

	return todo, nil
}
