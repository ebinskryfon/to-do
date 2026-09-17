package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"todo/internal/domain/entity"
	domainerrors "todo/internal/domain/errors"
	"todo/internal/domain/repository"
)

type todoModel struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Title       string     `gorm:"column:title;type:varchar(255);not null"`
	Description string     `gorm:"column:description;type:text"`
	Completed   bool       `gorm:"column:completed;not null;default:false"`
	IsActive    bool       `gorm:"column:is_active;not null;default:true"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;not null"`
	DeletedAt   *time.Time `gorm:"column:deleted_at"`
	CreatedBy   *uuid.UUID `gorm:"column:created_by"`
	UpdatedBy   *uuid.UUID `gorm:"column:updated_by"`
	DeletedBy   *uuid.UUID `gorm:"column:deleted_by"`
}

func (todoModel) TableName() string {
	return "todos"
}

func (m *todoModel) toEntity() entity.Todo {
	return entity.Todo{
		SoftDeletableEntity: entity.SoftDeletableEntity{
			AuditableEntity: entity.AuditableEntity{
				BaseModel: entity.BaseModel{
					CreatedAt: m.CreatedAt,
					UpdatedAt: m.UpdatedAt,
					DeletedAt: m.DeletedAt,
				},
				CreatedBy: m.CreatedBy,
				UpdatedBy: m.UpdatedBy,
			},
			IsActive:  m.IsActive,
			DeletedBy: m.DeletedBy,
		},
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		Completed:   m.Completed,
	}
}

func fromEntity(t *entity.Todo) todoModel {
	return todoModel{
		ID:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		Completed:   t.Completed,
		IsActive:    t.IsActive,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		DeletedAt:   t.DeletedAt,
		CreatedBy:   t.CreatedBy,
		UpdatedBy:   t.UpdatedBy,
		DeletedBy:   t.DeletedBy,
	}
}

type todoRepository struct {
	db *gorm.DB
}

var _ repository.TodoRepository = (*todoRepository)(nil)

// NewTodoRepository creates a new PostgreSQL GORM implementation of repository.TodoRepository.
func NewTodoRepository(db *gorm.DB) repository.TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(ctx context.Context, todo *entity.Todo) error {
	m := fromEntity(todo)
	if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
		return err
	}
	*todo = m.toEntity()
	return nil
}

func (r *todoRepository) GetByID(ctx context.Context, id uuid.UUID) (*entity.Todo, error) {
	var m todoModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND is_active = ? AND deleted_at IS NULL", id, true).
		First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainerrors.ErrTodoNotFound
		}
		return nil, err
	}
	t := m.toEntity()
	return &t, nil
}

func (r *todoRepository) Update(ctx context.Context, todo *entity.Todo) error {
	result := r.db.WithContext(ctx).
		Model(&todoModel{}).
		Where("id = ? AND is_active = ? AND deleted_at IS NULL", todo.ID, true).
		Updates(map[string]interface{}{
			"title":       todo.Title,
			"description": todo.Description,
			"completed": todo.Completed,
			"updated_at":  todo.UpdatedAt,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domainerrors.ErrTodoNotFound
	}

	return nil
}
