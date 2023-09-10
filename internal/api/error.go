package api

import "net/http"

type RequestError struct {
	Err         error
	UserMessage *string
	StatusCode  int
}

func NewError(statusCode int, err error) *RequestError {
	return &RequestError{
		Err:        err,
		StatusCode: statusCode,
	}
}

func NewInternalServerError(err error, message string) *RequestError {
	return NewErrorWithMessage(http.StatusInternalServerError, err, message)
}

func NewDatabaseError(err error) *RequestError {
	return NewInternalServerError(err, "Error interacting with database: please contact support")
}

func NewErrorWithMessage(statusCode int, err error, message string) *RequestError {
	return &RequestError{
		Err:         err,
		UserMessage: &message,
		StatusCode:  statusCode,
	}
}

func (r *RequestError) Error() string {
	if r.UserMessage != nil {
		return *r.UserMessage
	} else {
		return r.Err.Error()
	}
}
