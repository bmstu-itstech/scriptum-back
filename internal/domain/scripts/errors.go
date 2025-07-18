package scripts

import (
	"errors"
)

var (
	ErrFieldNameEmpty = errors.New("field name cannot be empty")
	ErrFieldNameLen   = errors.New("field name cannot exceed given length")

	ErrFieldDescEmpty = errors.New("field description cannot be empty")
	ErrFieldDescLen   = errors.New("field description cannot exceed given length")

	ErrFieldUnitEmpty = errors.New("field unit cannot be empty")
	ErrFieldUnitLen   = errors.New("field unit cannot exceed given length")

	ErrFileNameEmpty    = errors.New("file name cannot be empty")
	ErrFileTypeEmpty    = errors.New("file type cannot be empty")
	ErrFileContentEmpty = errors.New("file content cannot be empty")
	ErrVectorEmpty      = errors.New("vector cannot be empty")
	ErrFieldsEmpty      = errors.New("fields cannot be empty")

	ErrPathEmpty = errors.New("script path cannot be empty")
	ErrPathLen   = errors.New("script path cannot exceed given length")

	ErrFullNameEmpty = errors.New("full name cannot be empty")
	ErrEmailEmpty    = errors.New("email cannot be empty")

	ErrNameEmpty = errors.New("script name cannot be empty")
	ErrNameLen   = errors.New("script name cannot exceed given length")

	ErrDescriptionEmpty = errors.New("script description cannot be empty")
	ErrDescriptionLen   = errors.New("script description cannot exceed given length")

	ErrInvalidJobID       = errors.New("invalid job id")
	ErrInvalidUserID      = errors.New("invalid user id")
	ErrInvalidCommand     = errors.New("invalid command")
	ErrInvalidInterpreter = errors.New("invalid interpreter")
	ErrInvalidType        = errors.New("invalid type")
	ErrInvalidVisibility  = errors.New("invalid visibility")
	ErrInvalidValueType   = errors.New("invalid value type")
	ErrJobNotExists       = errors.New("job not exists")
	ErrComplexConversion  = errors.New("cannot convert to complex")
	ErrRealConversion     = errors.New("cannot convert to real")
	ErrIntegerConversion  = errors.New("cannot convert to integer")

	ErrInvalidSessionService  = errors.New("sessionService is nil")
	ErrInvalidScriptService   = errors.New("scriptService is nil")
	ErrInvalidUserService     = errors.New("userService is nil")
	ErrInvalidResService      = errors.New("resService is nil")
	ErrInvalidLauncherService = errors.New("launcherService is nil")
	ErrInvalidJobService      = errors.New("jobService is nil")
	ErrInvalidNotifierService = errors.New("notifierService is nil")
	ErrInvalidManagerService  = errors.New("managerService is nil")

	ErrNotAdmin         = errors.New("not admin")
	ErrNoAccessToDelete = errors.New("user has no access to delete script")
	ErrNoAccessToUpdate = errors.New("user has no access to update script")

	ErrFieldCnt = errors.New("число значений вывода не совпадает с числом полей ")

	ErrScriptLaunch = errors.New("ошибка запуска скрипта")
)
