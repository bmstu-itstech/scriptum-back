package usecase

// import (
// 	"context"

// 	authpb "github.com/bmstu-itstech/scriptum-back/auth"
// 	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
// )

// type UserLoginUC struct {
// 	authS authpb.SessionServiceClient
// }

// func NewUserLoginUC(authS authpb.SessionServiceClient) (*UserLoginUC, error) {
// 	if authS == nil {
// 		return nil, scripts.ErrInvalidSessionService
// 	}
// 	return &UserLoginUC{authS: authS}, nil
// }

// func (l *UserLoginUC) Login(ctx context.Context, login, password string) (string, error) {
// 	resp, err := l.authS.Login(ctx, &authpb.LoginRequest{Login: login, Password: password})
// 	if err != nil {
// 		return "", err
// 	}
// 	return resp.Token(), nil
// }
