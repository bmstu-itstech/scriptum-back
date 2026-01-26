package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zhikh23/pgutils"

	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

func (r *Repository) SaveBlueprint(ctx context.Context, blueprint *entity.Blueprint) error {
	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		rB := blueprintRowFromDomain(blueprint)
		if err := r.insertBlueprintRow(ctx, tx, rB); err != nil {
			return err
		}
		if len(blueprint.In()) > 0 {
			rIn := blueprintFieldRowsFromDomain(blueprint.In(), blueprint.ID())
			if err := r.insertBlueprintInputFieldRows(ctx, tx, rIn); err != nil {
				return err
			}
		}
		if len(blueprint.Out()) > 0 {
			rOut := blueprintFieldRowsFromDomain(blueprint.Out(), blueprint.ID())
			if err := r.insertBlueprintOutputFieldRows(ctx, tx, rOut); err != nil {
				return err
			}
		}
		return nil
	})
	if pgutils.IsUniqueViolationError(err) {
		return fmt.Errorf("%w: %s", ports.ErrJobAlreadyExists, string(blueprint.ID()))
	}
	return err
}

func (r *Repository) DeleteBlueprint(ctx context.Context, id value.BlueprintID) error {
	return pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		err := r.softDeleteBlueprintRow(ctx, r.db, string(id))
		if errors.Is(err, pgutils.ErrNoAffectedRows) {
			return fmt.Errorf("%w: %s", ports.ErrBlueprintNotFound, id)
		}
		return r.softDeleteBlueprintJobRows(ctx, r.db, string(id))
	})
}
