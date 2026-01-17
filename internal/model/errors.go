package model

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidInput  = errors.New("invalid input")
	ErrUnauthorized  = errors.New("unauthorized")
)

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on %s: %s", e.Field, e.Message)
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{Field: field, Message: message}
}

type OperationError struct {
	Op     string
	Entity string
	ID     string
	Err    error
}

func (e *OperationError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("%s %s %s: %v", e.Op, e.Entity, e.ID, e.Err)
	}
	return fmt.Sprintf("%s %s: %v", e.Op, e.Entity, e.Err)
}

func (e *OperationError) Unwrap() error {
	return e.Err
}

func NewOperationError(op, entity, id string, err error) *OperationError {
	return &OperationError{Op: op, Entity: entity, ID: id, Err: err}
}
