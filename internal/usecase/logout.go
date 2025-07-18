package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type UserLogoutUC struct {
	authS scripts.AuthService
}

func NewUserLogoutUC(authS scripts.AuthService) (*UserLogoutUC, error) {
	if authS == nil {
		return nil, scripts.ErrInvalidSessionService
	}
	return &UserLogoutUC{authS: authS}, nil
}

func (u *UserLogoutUC) Logout(ctx context.Context, token string) error {
	err := u.authS.Logout(ctx, token)
	return err
}
