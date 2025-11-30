package auth

import "context"

type ctxKey int

const ctxUserIDKey ctxKey = iota

func ExtractUserIDFromContext(ctx context.Context) (int64, bool) {
	uid, ok := ctx.Value(ctxUserIDKey).(int64)
	return uid, ok
}

func contextWithUserID(ctx context.Context, uid int64) context.Context {
	return context.WithValue(ctx, ctxUserIDKey, uid)
}
