package validation

import "fmt"

type InvalidInputError struct {
	Message string
}

func (e *InvalidInputError) Error() string {
	return e.Message
}

func NewInvalidInputError(message string) *InvalidInputError {
	return &InvalidInputError{Message: message}
}

func NewInvalidInputErrorf(message string, args ...any) *InvalidInputError {
	return &InvalidInputError{Message: fmt.Sprintf(message, args...)}
}
