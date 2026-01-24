package ports

import (
	"context"
	"errors"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

var ErrTokenInvalid = errors.New("token invalid")

type TokenService interface {
	GenerateToken(ctx context.Context, userID value.UserID) (value.Token, error)
	VerifyToken(ctx context.Context, token value.Token) (value.UserID, error)
}
