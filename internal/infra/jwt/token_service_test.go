package jwt_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/config"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/jwt"
)

func TestTokenService(t *testing.T) {
	cfg := config.JWT{
		Secret:    "a-very-very-secret-string",
		AccessTTL: time.Second,
	}
	service, err := jwt.NewTokenService(cfg)
	require.NoError(t, err)
	uid := value.NewUserID()

	token, err := service.GenerateToken(context.Background(), uid)
	require.NoError(t, err)

	parsedUID, err := service.VerifyToken(context.Background(), token)
	require.NoError(t, err)
	require.Equal(t, uid, parsedUID)

	require.Eventually(t, func() bool {
		_, err = service.VerifyToken(context.Background(), token)
		return errors.Is(err, ports.ErrTokenInvalid) // Срок годности токена истёк
	}, 2*time.Second, time.Second)
}
