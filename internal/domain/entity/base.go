package entity

import (
	"time"

	"github.com/google/uuid"
)

// BaseModel is the minimum every persisted entity carries.
type BaseModel struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

// AuditableEntity tracks who created/updated the row.
type AuditableEntity struct {
	BaseModel
	CreatedBy *uuid.UUID `json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID `json:"updated_by,omitempty"`
}

// SoftDeletableEntity adds an active flag and records who soft-deleted.
type SoftDeletableEntity struct {
	AuditableEntity
	IsActive  bool       `json:"is_active"`
	DeletedBy *uuid.UUID `json:"deleted_by,omitempty"`
}
