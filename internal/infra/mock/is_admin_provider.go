package mock

import (
	"context"
	"log"

	ssov1 "github.com/BOBAvov/protos_sso/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/bmstu-itstech/scriptum-back/internal/config"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type SSO struct {
	api ssov1.AuthClient
}

func NewSSOClient(config config.SSO) *SSO {
	addr := config.Host + ":" + config.Port

	cc, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("did not connect: ", err)
		return nil
	}
	return &SSO{
		api: ssov1.NewAuthClient(cc),
	}
}

func (i *SSO) IsAdmin(ctx context.Context, uid value.UserID) (bool, error) {
	resp, err := i.api.IsAdmin(ctx, &ssov1.IsAdminRequest{
		UserId: int64(uid),
	})
	if err != nil {
		return false, err
	}
	return resp.IsAdmin, nil
}
