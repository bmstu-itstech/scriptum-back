package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type UserDeleteUC struct {
	userS scripts.UserRepository
}

func NewUserDeleteUC(userS scripts.UserRepository) (*UserDeleteUC, error) {
	if userS == nil {
		return nil, scripts.ErrInvalidUserService
	}
	return &UserDeleteUC{userS: userS}, nil
}

func (u *UserDeleteUC) DeleteUser(ctx context.Context, actorID, userID uint32) error {
	maybeAdmin, err := u.userS.User(ctx, scripts.UserID(actorID))
	if err != nil {
		return err
	}
	if !maybeAdmin.IsAdmin() {
		return scripts.ErrNotAdmin
	}

	return u.userS.DeleteUser(ctx, scripts.UserID(userID))
}
