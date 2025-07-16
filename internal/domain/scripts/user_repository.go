package scripts

import "context"

type UserRepository interface {
	User(ctx context.Context, id UserID) (*User, error)
}
