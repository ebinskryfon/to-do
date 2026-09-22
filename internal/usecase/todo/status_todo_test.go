package todo_test

import (
	"context"
	"testing"

	"github.com/google/uuid"

	"todo/internal/domain/entity"
	"todo/internal/domain/repository"
	"todo/internal/usecase/todo"
)

type mockUpdateStatusTodoRepo struct {
	gotID        uuid.UUID
	gotCompleted bool
	statusErr    error
}

var _ repository.TodoRepository = (*mockUpdateStatusTodoRepo)(nil)

func (m *mockUpdateStatusTodoRepo) Create(ctx context.Context, t *entity.Todo) error {
	return nil
}

func (m *mockUpdateStatusTodoRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Todo, error) {
	return nil, nil
}

func (m *mockUpdateStatusTodoRepo) Update(ctx context.Context, t *entity.Todo) error {
	return nil
}

func (m *mockUpdateStatusTodoRepo) List(ctx context.Context, page, pageSize int, completed *bool) ([]entity.Todo, int, error) {
	return nil, 0, nil
}

func (m *mockUpdateStatusTodoRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockUpdateStatusTodoRepo) UpdateStatus(ctx context.Context, id uuid.UUID, completed bool) error {
	m.gotID = id
	m.gotCompleted = completed
	return m.statusErr
}

func TestUpdateTodoStatusUsecase_SuccessTrue(t *testing.T) {
	targetID := uuid.New()

	repo := &mockUpdateStatusTodoRepo{}
	uc := todo.NewUpdateTodoStatusUsecase(repo)

	err := uc.Execute(context.Background(), targetID, true)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.gotID != targetID {
		t.Fatalf("expected ID %v, got %v", targetID, repo.gotID)
	}

	if repo.gotCompleted != true {
		t.Fatalf("expected completed true, got %v", repo.gotCompleted)
	}
}

func TestUpdateTodoStatusUsecase_SuccessFalse(t *testing.T) {
	targetID := uuid.New()

	repo := &mockUpdateStatusTodoRepo{}
	uc := todo.NewUpdateTodoStatusUsecase(repo)

	err := uc.Execute(context.Background(), targetID, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.gotID != targetID {
		t.Fatalf("expected ID %v, got %v", targetID, repo.gotID)
	}

	if repo.gotCompleted != false {
		t.Fatalf("expected completed false, got %v", repo.gotCompleted)
	}
}

func TestUpdateTodoStatusUsecase_RepoErr(t *testing.T) {
	targetID := uuid.New()

	expectedErr := context.Canceled

	repo := &mockUpdateStatusTodoRepo{
		statusErr: expectedErr,
	}

	uc := todo.NewUpdateTodoStatusUsecase(repo)

	err := uc.Execute(context.Background(), targetID, true)
	if err != expectedErr {
		t.Fatalf("expected %v, got %v", expectedErr, err)
	}
}