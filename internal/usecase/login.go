package usecase

// import (
// 	"context"

// 	authpb "github.com/bmstu-itstech/scriptum-back/auth"
// 	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
// )

// type UserLoginUC interface {
// 	Login(ctx context.Context, login, password string) (string, error)
// }

// type UserLoginUCImp struct {
// 	sessionService authpb.SessionServiceClient
// }

// func NewUserLoginUCImp(sessionService authpb.SessionServiceClient) (*UserLoginUCImp, error) {
// 	if sessionService == nil {
// 		return nil, scripts.ErrInvalidSessionService
// 	}
// 	return &UserLoginUCImp{sessionService: sessionService}, nil
// }

// func (l *UserLoginUCImp) Login(ctx context.Context, login, password string) (string, error) {
// 	resp, err := l.sessionService.Login(ctx, &authpb.LoginRequest{Login: login, Password: password})
// 	if err != nil {
// 		return "", err
// 	}
// 	return resp.Token(), nil
// }
