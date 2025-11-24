package ports

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type BoxRepository interface {
	BoxProvider

	SaveBox(ctx context.Context, box *entity.Box) error
	DeleteBox(ctx context.Context, id value.BoxID) error
}
