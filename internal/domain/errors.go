package domain

import (
	"errors"
	"fmt"
)

var ErrPermissionDenied = errors.New("permission denied")

type InvalidInputError struct {
	Code    string
	Message string
}

func NewInvalidInputError(code string, message string) InvalidInputError {
	return InvalidInputError{
		Code:    code,
		Message: message,
	}
}

func (e InvalidInputError) Error() string {
	return fmt.Sprintf("invalid input: %s", e.Code)
}
