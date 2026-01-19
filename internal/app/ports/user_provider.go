package ports

import (
	"context"
	"errors"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

var ErrUserNotFound = errors.New("user is not found")

type UserProvider interface {
	User(ctx context.Context, id value.UserID) (*entity.User, error)
	UserByEmail(ctx context.Context, email value.Email) (*entity.User, error)
}
