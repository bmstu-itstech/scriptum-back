package jwtauth

import (
	"context"
	"errors"
)

type ctxKey int

const userCtxKey ctxKey = iota

func userIDToContext(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userCtxKey, userID)
}

var ErrNoUserIDInContext = errors.New("no user id in context")

func UserIDFromContext(ctx context.Context) (int64, error) {
	u, ok := ctx.Value(userCtxKey).(int64)
	if !ok {
		return 0, ErrNoUserIDInContext
	}
	return u, nil
}
