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
		panic(scripts.ErrInvalidUserRepository)
	}
	return &GetUserUC{userR: userR}, nil
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
