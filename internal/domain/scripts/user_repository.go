package scripts

import "context"

type UserRepository interface {
	User(ctx context.Context, id uint32) (*User, error)
}
