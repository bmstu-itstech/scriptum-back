package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type UserUpdateUC struct {
	userS scripts.UserRepository
}

func NewUserUpdateUC(userS scripts.UserRepository) (*UserUpdateUC, error) {
	if userS == nil {
		return nil, scripts.ErrInvalidUserService
	}
	return &UserUpdateUC{userS: userS}, nil
}

func (u *UserUpdateUC) UpdateUser(ctx context.Context, actorID uint32, dto UserDTO) error {
	maybeAdmin, err := u.userS.User(ctx, scripts.UserID(actorID))
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
	err = u.userS.UpdateUser(ctx, user)
	return err
}
