package sso

import (
	"context"
	ssov1 "github.com/BOBAvov/protos_sso/gen/go/sso"
	"github.com/bmstu-itstech/scriptum-back/internal/config"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log/slog"
	"net"
	"strconv"
)

type SSO struct {
	conn  *grpc.ClientConn
	api   ssov1.AuthClient
	l     *slog.Logger
	AppId int32
}

func MustNewSSOClient(config config.SSO, l *slog.Logger) (*SSO, func() error) {
	sso, closeFn, err := NewSSOClient(config, l)
	if err != nil {
		panic("Error creating SSO client: " + err.Error())
	}
	return sso, closeFn
}

func NewSSOClient(config config.SSO, l *slog.Logger) (*SSO, func() error, error) {
	addr := net.JoinHostPort(config.Host, config.Port)
	cc, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, func() error { return nil }, err
	}

	closeFn := cc.Close

	return &SSO{
		api:   ssov1.NewAuthClient(cc),
		conn:  cc,
		l:     l,
		AppId: config.AppId,
	}, closeFn, nil
}

// IsAdmin checks if the user with the given uid has admin privileges. uid is int64!!!
func (s *SSO) IsAdmin(ctx context.Context, uid value.UserID) (bool, error) {
	const op = "infra.SSO.IsAdmin"

	l := s.l.With(
		slog.String("op", op),
		slog.Int64("uid", int64(uid)),
	)

	l.Debug("Checking admin status")
	// Добавляем метаданные с AppId
	ctx = metadata.NewOutgoingContext(ctx, metadata.Pairs(
		"x-user-id", strconv.FormatInt(int64(s.AppId), 10),
	))
	// Вызываем метод IsAdmin на сервере SSO
	resp, err := s.api.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: int64(uid),
	})
	if err != nil {
		l.Error("Failed to check admin status: ", err.Error())
		return false, err
	}

	l.Debug("Admin status: ", resp.IsAdmin)

	return resp.IsAdmin, nil
}
