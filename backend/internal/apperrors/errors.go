package apperrors

import "fmt"

var (
	ErrNotFound     = fmt.Errorf("not found")
	ErrConflict     = fmt.Errorf("conflict")
	ErrInvalidInput = fmt.Errorf("invalid input")
)

func NewInvalidInput(message string) error {
	return fmt.Errorf("%w: %s", ErrInvalidInput, message)
}

func NewConflict(message string) error {
	return fmt.Errorf("%w: %s", ErrConflict, message)
}
