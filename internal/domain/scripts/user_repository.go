package scripts

import "context"

type UserRepository interface {
	User(ctx context.Context, id UserID) (*User, error)
	Users(ctx context.Context) ([]User, error)
	StoreUser(ctx context.Context, user User) (UserID, error)
	DeleteUser(ctx context.Context, userID UserID) error
	UpdateUser(ctx context.Context, user User) error
}
