package server

import (
	"context"
	"log/slog"
	"net"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/bmstu-itstech/scriptum-back/pkg/server/auth"
)

func RunGRPCServerOnAddr(ctx context.Context, l *slog.Logger, addr string, register func(server *grpc.Server)) error {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.StartCall, logging.FinishCall,
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) error {
			l.Error("Recovered from panic", slog.Any("panic", p))
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(interceptorLogger(l), loggingOpts...),
		auth.UnaryServerInterceptor(),
	))
	register(s)
	reflection.Register(s)

	go func() {
		l.InfoContext(ctx, "gRPC server starts listening", slog.String("addr", addr))
		if err = s.Serve(lis); err != nil {
			l.ErrorContext(ctx, "gRPC server failed to serve", slog.String("error", err.Error()))
		}
	}()

	<-ctx.Done()

	l.InfoContext(ctx, "gRPC server shutting down", slog.String("addr", addr))
	s.GracefulStop()
	return nil
}

func interceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
