package jwtauth

import (
	"context"
	"errors"
)

type ctxKey int

const userCtxKey ctxKey = iota

func userIDToContext(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userCtxKey, userID)
}

var ErrNoUserIDInContext = errors.New("no user id in context")

func UserIDFromContext(ctx context.Context) (string, error) {
	u, ok := ctx.Value(userCtxKey).(string)
	if !ok {
		return "", ErrNoUserIDInContext
	}
	return u, nil
}
