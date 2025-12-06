package auth

import (
	"context"
	"strconv"

	"google.golang.org/grpc/metadata"
)

func ClientOutgoingContext(ctx context.Context, uid int64) context.Context {
	md := metadata.New(map[string]string{"x-user-id": strconv.FormatInt(uid, 10)})
	return metadata.NewOutgoingContext(ctx, md)
}
