package scripts

import (
	"errors"
)

var ErrInvalidInput = errors.New("invalid input")
var ErrInvalidStateChange = errors.New("invalid state change")
var ErrPermissionDenied = errors.New("permission denied")
