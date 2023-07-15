package validation

type InvalidInputError struct {
	Message string
}

func (e *InvalidInputError) Error() string {
	return e.Message
}

func NewInvalidInputError(message string) *InvalidInputError {
	return &InvalidInputError{Message: message}
}
