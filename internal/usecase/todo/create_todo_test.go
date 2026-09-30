package todo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"todo/internal/domain/entity"
	domainerrors "todo/internal/domain/errors"
	"todo/internal/domain/repository"
	"todo/internal/usecase/todo"
)

type mockTodoRepo struct {
	createdTodo *entity.Todo
	createErr   error
}

var _ repository.TodoRepository = (*mockTodoRepo)(nil)

func (m *mockTodoRepo) Create(ctx context.Context, t *entity.Todo) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.createdTodo = t
	return nil
}

func (m *mockTodoRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Todo, error) {
	return nil, nil
}

func (m *mockTodoRepo) List(ctx context.Context, page, pageSize int, completed *bool) ([]entity.Todo, int, error) {
	return nil, 0, nil
}

func (m *mockTodoRepo) Update(ctx context.Context, t *entity.Todo) error {
	return nil
}

func (m *mockTodoRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockTodoRepo) UpdateStatus(ctx context.Context, id uuid.UUID, completed bool) error {
	return nil
}

func TestCreateTodoUsecase_Success(t *testing.T) {
	repo := &mockTodoRepo{}
	uc := todo.NewCreateTodoUsecase(repo)

	created, err := uc.Execute(context.Background(), "Test Title", "Test Description")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if created == nil {
		t.Fatal("expected created todo, got nil")
	}

	if created.Title != "Test Title" {
		t.Errorf("expected title 'Test Title', got '%s'", created.Title)
	}

	if created.Description != "Test Description" {
		t.Errorf("expected description 'Test Description', got '%s'", created.Description)
	}

	if created.Completed {
		t.Errorf("expected completed to be false")
	}

	if !created.IsActive {
		t.Errorf("expected is_active to be true")
	}
}

func TestCreateTodoUsecase_EmptyTitle(t *testing.T) {
	repo := &mockTodoRepo{}
	uc := todo.NewCreateTodoUsecase(repo)

	_, err := uc.Execute(context.Background(), "   ", "Test Description")
	if err != domainerrors.ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}
