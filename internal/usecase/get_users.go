package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetUsersUC struct {
	userS scripts.UserRepository
}

func NewGetUsersUC(userS scripts.UserRepository) (*GetUsersUC, error) {
	if userS == nil {
		return nil, scripts.ErrInvalidUserService
	}

	return &GetUsersUC{userS: userS}, nil
}

func (u *GetUsersUC) GetUsers(ctx context.Context) ([]UserDTO, error) {
	users, err := u.userS.Users(ctx)
	if err != nil {
		return nil, err
	}

	dto := make([]UserDTO, 0, len(users))
	for _, u := range users {
		dto = append(dto, UserToDTO(u))
	}

	return dto, nil
}
