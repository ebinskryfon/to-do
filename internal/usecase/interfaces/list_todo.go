package interfaces

import (
	"context"
	"todo/internal/domain/entity"
)

type ListTodoUsecase interface {
	Execute(ctx context.Context, page, pageSize int, completed *bool) ([]entity.Todo, int, error)
}