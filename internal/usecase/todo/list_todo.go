package todo

import (
	"context"
	"todo/internal/domain/entity"
	"todo/internal/domain/repository"
	"todo/internal/usecase/interfaces"
)

type listTodoUsecase struct {
	repo repository.TodoRepository
}

func NewListTodoUsecase(repo repository.TodoRepository) interfaces.ListTodoUsecase {
	return &listTodoUsecase{repo: repo}
}

func (uc *listTodoUsecase) Execute(ctx context.Context, page, pageSize int, completed *bool) ([]entity.Todo, int, error) {
	return uc.repo.List(ctx, page, pageSize, completed)
}