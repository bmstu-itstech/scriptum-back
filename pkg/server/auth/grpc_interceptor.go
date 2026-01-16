package auth

import (
	"context"
	"strconv"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if uidHeaders := md.Get("x-user-id"); len(uidHeaders) > 0 {
				if uid, err := strconv.ParseInt(uidHeaders[0], 10, 64); err == nil {
					ctx = contextWithUserID(ctx, uid)
				}
			}
		}
		return handler(ctx, req)
	}
}
