package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetUserUC struct {
	userS scripts.UserRepository
}

func NewGetUserUC(userS scripts.UserRepository) (*GetUserUC, error) {
	if userS == nil {
		return nil, scripts.ErrInvalidUserService
	}
	return &GetUserUC{userS: userS}, nil
}

func (u *GetUserUC) GetUser(ctx context.Context, userID uint32) (UserDTO, error) {
	user, err := u.userS.User(ctx, scripts.UserID(userID))
	if err != nil {
		return UserDTO{}, err
	}
	return UserToDTO(*user), nil
}
