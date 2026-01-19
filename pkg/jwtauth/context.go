package jwtauth

import (
	"context"
)

type ctxKey int

const ctxKeyUID ctxKey = iota

func toContext(ctx context.Context, uid string) context.Context {
	return context.WithValue(ctx, ctxKeyUID, uid)
}

func FromContext(ctx context.Context) (string, bool) {
	uid, ok := ctx.Value(ctxKeyUID).(string)
	return uid, ok
}
