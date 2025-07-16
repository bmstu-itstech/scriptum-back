package usecase

import (
	"context"

	authpb "github.com/bmstu-itstech/scriptum-back/auth"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type UserLogoutUC struct {
	authS authpb.SessionServiceClient
}

func NewUserLogoutUC(authS authpb.SessionServiceClient) (*UserLogoutUC, error) {
	if authS == nil {
		return nil, scripts.ErrInvalidSessionService
	}
	return &UserLogoutUC{authS: authS}, nil
}

func (u *UserLogoutUC) Logout(ctx context.Context, token string) error {
	_, err := u.authS.Logout(ctx, &authpb.LogoutRequest{Token: token})
	return err
}
