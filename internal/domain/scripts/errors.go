package scripts

import (
	"errors"
)

var (
	ErrFieldNameEmpty        = errors.New("field name cannot be empty")
	ErrFieldDescEmpty        = errors.New("field description cannot be empty")
	ErrFieldUnitEmpty        = errors.New("field unit cannot be empty")
	ErrFileNameEmpty         = errors.New("file name cannot be empty")
	ErrFileTypeEmpty         = errors.New("file type cannot be empty")
	ErrFileContentEmpty      = errors.New("file content cannot be empty")
	ErrInvalidJobID          = errors.New("invalid job id")
	ErrInvalidUserID         = errors.New("invalid user id")
	ErrEmptyVector           = errors.New("vector cannot be empty")
	ErrInvalidCommand        = errors.New("invalid command")
	ErrInvalidInterpreter    = errors.New("invalid interpreter")
	ErrFieldsEmpty           = errors.New("fields cannot be empty")
	ErrPathEmpty             = errors.New("path cannot be empty")
	ErrFullNameEmpty         = errors.New("full name cannot be empty")
	ErrEmailEmpty            = errors.New("email cannot be empty")
	ErrInvalidType           = errors.New("invalid type")
	ErrInvalidSessionService = errors.New("sessionService is nil")
	ErrInvalidScriptService  = errors.New("scriptService is nil")
	ErrInvalidUserService    = errors.New("userService is nil")
	ErrInvalidVisibility     = errors.New("invalid visibility")
	ErrJobNotExists          = errors.New("job не найден")
	ErrComplexConversion     = errors.New("не удалось привести число к комплексному виду")
	ErrRealConversion        = errors.New("не удалось привести число к дробному виду")
	ErrIntegerConversion     = errors.New("не удалось привести число к целочисленному виду")
)
