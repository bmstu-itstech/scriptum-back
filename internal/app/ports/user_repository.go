package ports

import (
	"context"
	"errors"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

var ErrUserAlreadyExists = errors.New("user already exists")

type UserRepository interface {
	UserProvider

	SaveUser(ctx context.Context, u *entity.User) error
	UpdateUser(ctx context.Context, uid value.UserID, updateFn func(inner context.Context, u *entity.User) error) error
	DeleteUser(ctx context.Context, uid value.UserID) error
}
