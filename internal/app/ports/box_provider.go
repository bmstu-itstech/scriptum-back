package ports

import (
	"context"
	"errors"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

var ErrBoxNotFound = errors.New("box not found")

type BoxProvider interface {
	Box(ctx context.Context, id value.BoxID) (*entity.Box, error)
	// Boxes возвращает все Box, доступные пользователю
	Boxes(ctx context.Context, uid value.UserID) ([]*entity.Box, error)
	// SearchBoxes осуществляет нечёткий поиск по коллекции доступных пользователю Box по имени
	// или его части
	SearchBoxes(ctx context.Context, uid value.UserID, name string) ([]*entity.Box, error)
}
