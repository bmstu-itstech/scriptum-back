package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type UserDeleteUC struct {
	userR  scripts.UserRepository
	logger *slog.Logger
}

func NewUserDeleteUC(userR scripts.UserRepository, logger *slog.Logger) UserDeleteUC {
	if userR == nil {
		panic(scripts.ErrInvalidUserRepository)
	}
	if logger == nil {
		panic(scripts.ErrInvalidUserRepository)
	}
	return UserDeleteUC{userR: userR, logger: logger}
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
