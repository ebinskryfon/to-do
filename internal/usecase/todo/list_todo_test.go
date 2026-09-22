package todo_test

import (
	"context"
	"errors"
	"testing"

	"todo/internal/domain/entity"
	"todo/internal/domain/repository"
	"todo/internal/usecase/todo"

	"github.com/google/uuid"
)

type mockListTodoRepo struct {
	todos []entity.Todo
	total int
	listErr error
	gotPage int
	gotPageSize int
	gotCompleted *bool
}

var _ repository.TodoRepository = (*mockListTodoRepo)(nil)

func (m *mockListTodoRepo) Create(ctx context.Context, t *entity.Todo) error {
	return nil
}

func (m *mockListTodoRepo) GetByID(ctx context.Context, id uuid.UUID) (*entity.Todo, error) {
	return nil, nil
}

func (m *mockListTodoRepo) List(ctx context.Context, page, pageSize int, completed *bool) ([]entity.Todo, int, error) {
	m.gotPage = page
	m.gotPageSize = pageSize
	m.gotCompleted = completed

	return m.todos, m.total, m.listErr
}

func (m *mockListTodoRepo) Update(ctx context.Context, t *entity.Todo) error {
	return nil
}

func (m *mockListTodoRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func TestListTodoUsecase_Success(t *testing.T) {
	todo1 := entity.Todo{
		ID:    uuid.New(),
		Title: "Todo 1",
	}

	todo2 := entity.Todo{
		ID:    uuid.New(),
		Title: "Todo 2",
	}

	repo := &mockListTodoRepo{
		todos: []entity.Todo{todo1, todo2},
		total: 2,
	}

	uc := todo.NewListTodoUsecase(repo)

	todos, total, err := uc.Execute(context.Background(), 1, 10, nil)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if total != 2 {
		t.Errorf("expected total 2, got %d", total)
	}

	if len(todos) != 2 {
		t.Fatalf("expected 2 todos, got %d", len(todos))
	}

	if todos[0].Title != "Todo 1" {
		t.Errorf("expected title 'Todo 1', got '%s'", todos[0].Title)
	}

	if todos[1].Title != "Todo 2" {
		t.Errorf("expected title 'Todo 2', got '%s'", todos[1].Title)
	}

	if repo.gotPage != 1 {
		t.Fatalf("expected 1 got %v", repo.gotPage)
	}

	if repo.gotPageSize != 10 {
		t.Fatalf("expected 10 got %v", repo.gotPageSize)
	}

	if repo.gotCompleted != nil {
		t.Fatalf("expected nil got %v", repo.gotCompleted)
	}
}

func TestListTodoUsecase_CompletedTrue(t *testing.T) {
	completed := true

	repo := &mockListTodoRepo{
		todos: []entity.Todo{
			{
				ID:        uuid.New(),
				Title:     "Completed Todo",
				Completed: true,
			},
		},
		total: 1,
	}

	uc := todo.NewListTodoUsecase(repo)

	_, _, err := uc.Execute(context.Background(), 1, 10, &completed)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.gotCompleted == nil {
		t.Fatal("expected completed filter, got nil")
	}

	if *repo.gotCompleted != true {
		t.Errorf("expected completed true, got %v", *repo.gotCompleted)
	}
}

func TestListTodoUsecase_CompletedFalse(t *testing.T) {
	completed := false

	repo := &mockListTodoRepo{
		todos: []entity.Todo{
			{
				ID:        uuid.New(),
				Title:     "Pending Todo",
				Completed: false,
			},
		},
		total: 1,
	}

	uc := todo.NewListTodoUsecase(repo)

	_, _, err := uc.Execute(context.Background(), 2, 5, &completed)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if repo.gotPage != 2 {
		t.Errorf("expected page 2, got %d", repo.gotPage)
	}

	if repo.gotPageSize != 5 {
		t.Errorf("expected page size 5, got %d", repo.gotPageSize)
	}

	if repo.gotCompleted == nil {
		t.Fatal("expected completed filter, got nil")
	}

	if *repo.gotCompleted != false {
		t.Errorf("expected completed false, got %v", *repo.gotCompleted)
	}
}

func TestListTodoUsecase_RepositoryError(t *testing.T) {
	expectedErr := errors.New("repository error")

	repo := &mockListTodoRepo{
		listErr: expectedErr,
	}

	uc := todo.NewListTodoUsecase(repo)

	todos, total, err := uc.Execute(context.Background(), 1, 10, nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("expected repository error, got %v", err)
	}

	if todos != nil {
		t.Errorf("expected nil todos, got %v", todos)
	}

	if total != 0 {
		t.Errorf("expected total 0, got %d", total)
	}
}