package error

import "net/http"

type CustomError struct {
	Status  int
	Message string
	Errors  []string
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewBadRequestError(message string, errors ...string) error {
	return &CustomError{
		Status:  http.StatusBadRequest,
		Message: message,
		Errors:  errors,
	}
}

func NewForbiddenError(message string, errors ...string) error {
	return &CustomError{
		Status:  http.StatusForbidden,
		Message: message,
		Errors:  errors,
	}
}

func NewNotFoundError(message string, errors ...string) error {
	return &CustomError{
		Status:  http.StatusNotFound,
		Message: message,
		Errors:  errors,
	}
}

func NewRequestTimeoutError(message string, errors ...string) error {
	return &CustomError{
		Status:  http.StatusRequestTimeout,
		Message: message,
		Errors:  errors,
	}
}

func NewTooManyRequestError(message string, errors ...string) error {
	return &CustomError{
		Status:  http.StatusTooManyRequests,
		Message: message,
		Errors:  errors,
	}
}

func NewInternalServerError(message string, errors ...string) error {
	return &CustomError{
		Status:  http.StatusInternalServerError,
		Message: message,
		Errors:  errors,
	}
}
