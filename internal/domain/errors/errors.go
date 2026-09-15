package errors

import "errors"

var (
	// ErrTodoNotFound is returned when the requested todo item is not found.
	ErrTodoNotFound = errors.New("todo not found")
	// ErrInvalidInput is returned when input validation fails in the business layer.
	ErrInvalidInput = errors.New("invalid input")
)
