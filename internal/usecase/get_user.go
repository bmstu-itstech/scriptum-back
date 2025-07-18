package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetUserUC struct {
	userR scripts.UserRepository
}

func NewGetUserUC(userR scripts.UserRepository) (*GetUserUC, error) {
	if userR == nil {
		return nil, scripts.ErrInvalidUserService
	}
	return &GetUserUC{userR: userR}, nil
}

func (u *GetUserUC) GetUser(ctx context.Context, userID uint32) (UserDTO, error) {
	user, err := u.userR.User(ctx, scripts.UserID(userID))
	if err != nil {
		return UserDTO{}, err
	}
	return UserToDTO(*user), nil
}
