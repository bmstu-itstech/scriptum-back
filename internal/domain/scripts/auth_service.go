package scripts

import "context"

type AuthService interface {
	Login(ctx context.Context, login, password string) (string, error)
	Logout(ctx context.Context, token string) error
	CheckSession(ctx context.Context, token string)
}
