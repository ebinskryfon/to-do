package todo_test

import (
	"context"
	"testing"

	"todo/internal/domain/entity"
	"todo/internal/domain/repository"
	"todo/internal/usecase/todo"
	domainerrors "todo/internal/domain/errors"

	"github.com/google/uuid"
)

type mockDeleteTodoRepo struct {
	gotID uuid.UUID
	deleteErr error
}

var _ repository.TodoRepository = (*mockDeleteTodoRepo)(nil)

func (m *mockDeleteTodoRepo) Create(ctx context.Context, t *entity.Todo) error {
	return nil
}

func (m *mockDeleteTodoRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Todo, error) {
	return nil, nil
}

func (m *mockDeleteTodoRepo) List(ctx context.Context, page, pageSize int, completed *bool) ([]entity.Todo, int, error) {
	return nil, 0, nil
}

func (m *mockDeleteTodoRepo) Update(ctx context.Context, t *entity.Todo) error {
	return nil
}

func (m *mockDeleteTodoRepo) Delete(ctx context.Context, id uuid.UUID) error {
	m.gotID = id
	if m.deleteErr != nil {
		return m.deleteErr
	}

	return nil
}

func (m *mockDeleteTodoRepo) UpdateStatus(ctx context.Context, id uuid.UUID, completed bool) error {
	return nil
}

func TestDeleteTodoUsecase_Success(t *testing.T) {
	targetID := uuid.New()

	repo := &mockDeleteTodoRepo{}
	uc := todo.NewDeleteTodoUsecase(repo)

	err := uc.Execute(context.Background(), targetID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.gotID != targetID {
		t.Fatalf("expected same ids, got different ids")
	}
}

func TestDeleteTodoUsecase_RepoErr(t *testing.T) {
	targetID := uuid.New()
	repo := &mockDeleteTodoRepo{deleteErr: domainerrors.ErrTodoNotFound}
	uc := todo.NewDeleteTodoUsecase(repo)

	err := uc.Execute(context.Background(), targetID)
	if err != domainerrors.ErrTodoNotFound {
		t.Fatalf("expected %v, got %v", domainerrors.ErrTodoNotFound, err)
	}
}

func TestDeleteTodoUsecase_NotFound(t *testing.T) {
	targetID := uuid.New()
	repo := &mockDeleteTodoRepo{deleteErr: domainerrors.ErrTodoNotFound}
	uc := todo.NewDeleteTodoUsecase(repo)

	err := uc.Execute(context.Background(), targetID)
	if err != domainerrors.ErrTodoNotFound {
		t.Fatalf("expected %v, got %v", domainerrors.ErrTodoNotFound, err)
	}
}