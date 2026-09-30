package todo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"todo/internal/domain/entity"
	domainerrors "todo/internal/domain/errors"
	"todo/internal/usecase/todo"
)

type mockGetTodoRepo struct {
	todo   *entity.Todo
	getErr error
}

func (m *mockGetTodoRepo) Create(ctx context.Context, t *entity.Todo) error {
	return nil
}

func (m *mockGetTodoRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Todo, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.todo, nil
}

func (m *mockGetTodoRepo) List(ctx context.Context, page, pageSize int, completed *bool) ([]entity.Todo, int, error) {
	return nil, 0, nil
}

func (m *mockGetTodoRepo) Update(ctx context.Context, t *entity.Todo) error {
	return nil
}

func (m *mockGetTodoRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockGetTodoRepo) UpdateStatus(ctx context.Context, id uuid.UUID, completed bool) error {
	return nil
}

func TestGetTodoUsecase_Success(t *testing.T) {
	targetID := uuid.New()
	expectedTodo := &entity.Todo{
		ID:          targetID,
		Title:       "Sample Todo",
		Description: "Sample Description",
	}

	repo := &mockGetTodoRepo{todo: expectedTodo}
	uc := todo.NewGetTodoUsecase(repo)

	result, err := uc.Execute(context.Background(), targetID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil || result.ID != targetID {
		t.Fatalf("expected todo with ID %v, got %v", targetID, result)
	}
}

func TestGetTodoUsecase_NotFound(t *testing.T) {
	targetID := uuid.New()
	repo := &mockGetTodoRepo{getErr: domainerrors.ErrTodoNotFound}
	uc := todo.NewGetTodoUsecase(repo)

	_, err := uc.Execute(context.Background(), targetID)
	if err != domainerrors.ErrTodoNotFound {
		t.Fatalf("expected ErrTodoNotFound, got %v", err)
	}
}
