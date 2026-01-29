package ports

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type BlueprintRepository interface {
	Blueprint(ctx context.Context, id value.BlueprintID) (*entity.Blueprint, error)
	SaveBlueprint(ctx context.Context, box *entity.Blueprint) error
	DeleteBlueprint(ctx context.Context, id value.BlueprintID) error
}
