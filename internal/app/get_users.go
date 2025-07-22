package app

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetUsersUC struct {
	userR scripts.UserRepository
}

func NewGetUsersUC(userR scripts.UserRepository) (*GetUsersUC, error) {
	if userR == nil {
		return nil, scripts.ErrInvalidUserRepository
	}

	return &GetUsersUC{userR: userR}, nil
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
