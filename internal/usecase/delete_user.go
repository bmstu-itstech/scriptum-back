package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type UserDeleteUC struct {
	userR scripts.UserRepository
}

func NewUserDeleteUC(userR scripts.UserRepository) (*UserDeleteUC, error) {
	if userR == nil {
		return nil, scripts.ErrInvalidUserService
	}
	return &UserDeleteUC{userR: userR}, nil
}

func (u *UserDeleteUC) DeleteUser(ctx context.Context, actorID, userID uint32) error {
	maybeAdmin, err := u.userR.User(ctx, scripts.UserID(actorID))
	if err != nil {
		return err
	}
	if !maybeAdmin.IsAdmin() {
		return scripts.ErrNotAdmin
	}

	return u.userR.DeleteUser(ctx, scripts.UserID(userID))
}
