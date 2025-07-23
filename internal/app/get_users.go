package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetUsersUC struct {
	userR  scripts.UserRepository
	logger *slog.Logger
}

func NewGetUsersUC(userR scripts.UserRepository, logger *slog.Logger) GetUsersUC {
	if userR == nil {
		panic(scripts.ErrInvalidUserRepository)
	}
	if logger == nil {
		panic(scripts.ErrInvalidLogger)
	}

	return GetUsersUC{userR: userR, logger: logger}
}

func (u *GetUsersUC) GetUsers(ctx context.Context) ([]UserDTO, error) {
	users, err := u.userR.Users(ctx)
	if err != nil {
		return nil, err
	}

	dto := make([]UserDTO, 0, len(users))
	for _, u := range users {
		dto = append(dto, UserToDTO(u))
	}

	return dto, nil
}
