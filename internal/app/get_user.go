package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetUserUC struct {
	userR  scripts.UserRepository
	logger *slog.Logger
}

func NewGetUserUC(userR scripts.UserRepository, logger *slog.Logger) GetUserUC {
	if userR == nil {
		panic(scripts.ErrInvalidUserRepository)
	}
	if logger == nil {
		panic(scripts.ErrInvalidLogger)
	}
	return GetUserUC{userR: userR, logger: logger}
}

func (u *GetUserUC) GetUser(ctx context.Context, actorID, userID uint32) (UserDTO, error) {
	maybeAdmin, err := u.userR.User(ctx, scripts.UserID(actorID))
	if err != nil {
		return UserDTO{}, err
	}
	adm := maybeAdmin.IsAdmin()
	if !adm && userID != actorID {
		return UserDTO{}, scripts.ErrNoAccessToGet
	}
	user, err := u.userR.User(ctx, scripts.UserID(userID))
	if err != nil {
		return UserDTO{}, err
	}
	return UserToDTO(*user), nil
}
