package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type UserCreateUC struct {
	userS scripts.UserRepository
}

func NewUserCreateUC(userS scripts.UserRepository) (*UserCreateUC, error) {
	if userS == nil {
		return nil, scripts.ErrInvalidUserService
	}
	return &UserCreateUC{userS: userS}, nil
}

func (u *UserCreateUC) CreateUser(ctx context.Context, actorID uint32, newUser UserDTO) (uint32, error) {
	maybeAdmin, err := u.userS.User(ctx, scripts.UserID(actorID))
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
	userID, err := u.userS.StoreUser(ctx, user)
	if err != nil {
		return 0, err
	}
	
	return uint32(userID), nil
}
