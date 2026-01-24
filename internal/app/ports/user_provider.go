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
	Users(ctx context.Context) ([]*entity.User, error)
	UserByEmail(ctx context.Context, email string) (*entity.User, error)
}
