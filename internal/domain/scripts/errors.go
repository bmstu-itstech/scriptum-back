package scripts

import (
	"errors"
)

var (
	ErrFieldInvalid  = errors.New("error while creating field variable")
	ErrFileInvalid   = errors.New("error while creating file variable")
	ErrJobInvalid    = errors.New("error while creating job variable")
	ErrScriptInvalid = errors.New("error while creating script variable")
	ErrTypeInvalid   = errors.New("error while creating type variable")
	ErrValueInvalid  = errors.New("error while creating value variable")
	ErrVectorInvalid = errors.New("error while creating vector variable")
	ErrUserInvalid   = errors.New("error while creating user variable")

	ErrInvalidJobID       = errors.New("invalid job id")
	ErrInvalidUserID      = errors.New("invalid user id")
	ErrInvalidCommand     = errors.New("invalid command")
	ErrInvalidInterpreter = errors.New("invalid interpreter")
	ErrInvalidVisibility  = errors.New("invalid visibility")
	ErrInvalidValueType   = errors.New("invalid value type")
	ErrJobNotExists       = errors.New("job not exists")
	ErrComplexConversion  = errors.New("cannot convert to complex")
	ErrRealConversion     = errors.New("cannot convert to real")
	ErrIntegerConversion  = errors.New("cannot convert to integer")

	ErrInvalidScriptRepository = errors.New("scriptService is nil")
	ErrInvalidUserRepository   = errors.New("userService is nil")
	ErrInvalidResultRepository = errors.New("resService is nil")
	ErrInvalidLauncherService  = errors.New("launcherService is nil")
	ErrInvalidJobRepository    = errors.New("jobService is nil")
	ErrInvalidNotifierService  = errors.New("notifierService is nil")
	ErrInvalidManagerService   = errors.New("managerService is nil")

	ErrNotAdmin         = errors.New("not admin")
	ErrNoAccessToDelete = errors.New("user has no access to delete script")
	ErrNoAccessToUpdate = errors.New("user has no access to update script")
	ErrNoAccessToGet    = errors.New("user has no access to get this user")

	ErrFieldCount = errors.New("number of values in output does not match number of fields")

	ErrScriptLaunch = errors.New("script launch error")
)
