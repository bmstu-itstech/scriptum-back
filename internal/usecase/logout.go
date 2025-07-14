package usecase

// import (
// 	"context"

// 	authpb "github.com/bmstu-itstech/scriptum-back/auth"
// 	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
// )

// type UserLogoutUC interface {
// 	Logout(ctx context.Context, token string) error
// }

// type UserLogoutUCImp struct {
// 	sessionService authpb.SessionServiceClient
// }

// func (l *UserLogoutUCImp) SessionService() authpb.SessionServiceClient {
// 	return l.sessionService
// }

// func NewUserLogoutUCImp(sessionService authpb.SessionServiceClient) (*UserLogoutUCImp, error) {
// 	if sessionService == nil {
// 		return nil, scripts.ErrInvalidSessionService
// 	}
// 	return &UserLogoutUCImp{sessionService: sessionService}, nil
// }

// func (u *UserLogoutUCImp) Logout(ctx context.Context, token string) error {
// 	_, err := u.SessionService().Logout(ctx, &authpb.LogoutRequest{Token: token})
// 	return err
// }
