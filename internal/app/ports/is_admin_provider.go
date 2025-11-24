package ports

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type IsAdminProvider interface {
	IsAdmin(ctx context.Context, uid value.UserID) (bool, error)
}
