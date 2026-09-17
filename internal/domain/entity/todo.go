package entity

import "github.com/google/uuid"

// Todo represents a todo item with full audit and soft-delete capabilities.
type Todo struct {
	ID uuid.UUID `json:"id"`
	SoftDeletableEntity
	Title       string `json:"title"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}
