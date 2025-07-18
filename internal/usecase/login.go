package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type UserLoginUC struct {
	authS scripts.AuthService
}

func NewUserLoginUC(authS scripts.AuthService) (*UserLoginUC, error) {
	if authS == nil {
		return nil, scripts.ErrInvalidSessionService
	}
	return &UserLoginUC{authS: authS}, nil
}

func (l *UserLoginUC) Login(ctx context.Context, login, password string) (string, error) {
	resp, err := l.authS.Login(ctx, login, password)
	if err != nil {
		return "", err
	}
	return resp, nil
}
