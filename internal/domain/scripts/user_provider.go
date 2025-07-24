package scripts

import (
	"context"
	"errors"
)

var ErrUserNotFound = errors.New("user not found")

type UserProvider interface {
	User(ctx context.Context, id UserID) (*User, error)
}
