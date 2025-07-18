package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type UserUpdateUC struct {
	userR scripts.UserRepository
}

func NewUserUpdateUC(userR scripts.UserRepository) (*UserUpdateUC, error) {
	if userR == nil {
		return nil, scripts.ErrInvalidUserService
	}
	return &UserUpdateUC{userR: userR}, nil
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
