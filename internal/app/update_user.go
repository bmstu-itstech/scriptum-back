package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type UserUpdateUC struct {
	userR  scripts.UserRepository
	logger *slog.Logger
}

func NewUserUpdateUC(userR scripts.UserRepository, logger *slog.Logger) UserUpdateUC {
	if userR == nil {
		panic(scripts.ErrInvalidUserRepository)
	}
	if logger == nil {
		panic(scripts.ErrInvalidLogger)
	}
	return UserUpdateUC{userR: userR, logger: logger}
}

func (u *UserUpdateUC) UpdateUser(ctx context.Context, actorID uint32, dto UserDTO) error {
	maybeAdmin, err := u.userR.User(ctx, scripts.UserID(actorID))
	if err != nil {
		return err
	}
	if !maybeAdmin.IsAdmin() {
		return scripts.ErrNotAdmin
	}
	user, err := DTOToUser(dto)
	if err != nil {
		return err
	}
	err = u.userR.UpdateUser(ctx, user)
	return err
}
