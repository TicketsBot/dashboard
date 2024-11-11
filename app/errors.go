package app

import "fmt"

type ApiError struct {
	InternalError   error
	ExternalMessage string
}

var _ error = (*ApiError)(nil)

func NewError(internalError error, externalMessage string) *ApiError {
	return &ApiError{
		InternalError:   internalError,
		ExternalMessage: externalMessage,
	}
}

func NewServerError(internalError error) *ApiError {
	return NewError(internalError, "An internal server error occurred")
}

func (e *ApiError) Error() string {
	return fmt.Sprintf("internal error: %v, external message: %s", e.InternalError, e.ExternalMessage)
}

func (e *ApiError) Unwrap() error {
	return e.InternalError
}
