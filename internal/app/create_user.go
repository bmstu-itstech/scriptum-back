package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type UserCreateUC struct {
	userR  scripts.UserRepository
	logger *slog.Logger
}

func NewUserCreateUC(userR scripts.UserRepository, logger *slog.Logger) UserCreateUC {
	if userR == nil {
		panic(scripts.ErrInvalidUserRepository)
	}
	if logger == nil {
		panic(scripts.ErrInvalidLogger)
	}
	return UserCreateUC{userR: userR, logger: logger}
}

func (u *UserCreateUC) CreateUser(ctx context.Context, actorID uint32, newUser UserDTO) (uint32, error) {
	maybeAdmin, err := u.userR.User(ctx, scripts.UserID(actorID))
	if err != nil {
		return 0, err
	}
	if !maybeAdmin.IsAdmin() {
		return 0, scripts.ErrNotAdmin
	}

	user, err := DTOToUser(newUser)
	if err != nil {
		return 0, err
	}
	userID, err := u.userR.StoreUser(ctx, user)
	if err != nil {
		return 0, err
	}

	return uint32(userID), nil
}
