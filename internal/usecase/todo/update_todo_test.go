package todo_test

import (
	"context"
	"testing"
	"time"

	// "time"

	"todo/internal/domain/entity"
	domainerrors "todo/internal/domain/errors"
	"todo/internal/usecase/todo"

	"github.com/google/uuid"

	// domainerrors "todo/internal/domain/errors"
	"todo/internal/domain/repository"
	// "todo/internal/usecase/todo"
	"todo/internal/types"
)

type mockUpdateTodoRepo struct {
	todo *entity.Todo
	getErr error
	updatedTodo *entity.Todo
	updateErr error
}

var _ repository.TodoRepository = (*mockUpdateTodoRepo)(nil)

func (m *mockUpdateTodoRepo) Create(ctx context.Context, t *entity.Todo) error {
	return nil
}

func (m *mockUpdateTodoRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Todo, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	return m.todo, nil
}

func (m *mockUpdateTodoRepo) List(ctx context.Context, page, pageSize int, completed *bool) ([]entity.Todo, int, error) {
	return nil, 0, nil
}

func (m *mockUpdateTodoRepo) Update(ctx context.Context, t *entity.Todo) error {
	if m.updateErr != nil {
		return m.updateErr
	}

	m.updatedTodo = t
	return nil
}

func (m *mockUpdateTodoRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func TestUpdateTodoUsecase_Success(t *testing.T) {
	targetID := uuid.New()
	oldTime := time.Now().UTC().Add(-time.Hour)

	existingTodo := &entity.Todo{
		ID: targetID,
		Title: "Old Title",
		Description: "Old Description",
		Completed: false,
		SoftDeletableEntity: entity.SoftDeletableEntity{
			AuditableEntity: entity.AuditableEntity{
				BaseModel: entity.BaseModel{
					CreatedAt: oldTime,
					UpdatedAt: oldTime,
				},
			},
		},
	}

	repo := &mockUpdateTodoRepo{todo: existingTodo}
	uc := todo.NewUpdateTodoUsecase(repo)

	req := types.UpdateTodoRequest {
		Title: "New Title",
		Description: "New Description",
		Completed: true,
	}

	result, err := uc.Execute(context.Background(), targetID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatalf("expected updated todo, got nil")
	}

	if result.Title != "New Title" {
		t.Fatalf("expected 'New Title' as Title, got %v", result.Title)
	}

	if result.Description != "New Description" {
		t.Fatalf("expected 'New Description' as Description, got %v", result.Description)
	}

	if !result.Completed {
		t.Fatalf("expected Completed to be true")
	}

	if repo.updatedTodo == nil {
		t.Fatalf("expected updated todo, got nil")
	}

	if repo.updatedTodo.Title != "New Title" {
		t.Fatalf("expected 'New Title' as Title, got %v", repo.updatedTodo.Title)
	}

	if repo.updatedTodo.Description != "New Description" {
		t.Fatalf("expected 'New Description' as Description, got %v", repo.updatedTodo.Description)
	}

	if !repo.updatedTodo.Completed {
		t.Fatalf("expected Completed to be true")
	}
	
	if result.CreatedAt != oldTime {
		t.Fatalf("expected %v, got %v", oldTime, result.CreatedAt)
	}

	if time.Time.Equal(result.UpdatedAt, oldTime) {
		t.Fatalf("expected %v, got %v", result.UpdatedAt, oldTime)
	}
}

func TestUpdateTodoUsecase_EmptyTitle(t *testing.T) {
	targetID := uuid.New()

	repo := &mockUpdateTodoRepo{}
	uc := todo.NewUpdateTodoUsecase(repo)

	req := types.UpdateTodoRequest {Title: "   ",}

	_, err := uc.Execute(context.Background(), targetID, req)
	if err != domainerrors.ErrInvalidInput {
		t.Fatalf("expected ErrInvalidInput, got %v", err)
	}
}

func TestUpdateTodoUsecase_NotFound(t *testing.T) {
	targetID := uuid.New()

	repo := &mockUpdateTodoRepo{getErr: domainerrors.ErrTodoNotFound}
	uc := todo.NewUpdateTodoUsecase(repo)

	req := types.UpdateTodoRequest{
		Title: "Updated Title",
		Description: "Updated Description",
	}

	_, err := uc.Execute(context.Background(), targetID, req)
	if err != domainerrors.ErrTodoNotFound {
		t.Fatalf("expected ErrTodoNotFound, got %v", err)
	}
}