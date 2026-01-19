package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/zhikh23/pgutils"

	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

func (r *Repository) SaveBlueprint(ctx context.Context, blueprint *entity.Blueprint) error {
	l := r.l.With(
		slog.String("op", "postgres.Repository.SaveBlueprint"),
		slog.String("blueprint_id", string(blueprint.ID())),
	)
	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		rB := blueprintRowFromDomain(blueprint)
		if err := r.insertBlueprintRow(ctx, tx, rB); err != nil {
			return err
		}
		rIn := blueprintFieldRowsFromDomain(blueprint.In(), blueprint.ID())
		if err := r.insertBlueprintInputFieldRows(ctx, tx, rIn); err != nil {
			return err
		}
		rOut := blueprintFieldRowsFromDomain(blueprint.Out(), blueprint.ID())
		if err := r.insertBlueprintOutputFieldRows(ctx, tx, rOut); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		l.ErrorContext(ctx, "failed to execute transaction", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (r *Repository) DeleteBlueprint(ctx context.Context, id value.BlueprintID) error {
	l := r.l.With(
		slog.String("op", "postgres.Repository.DeleteBlueprint"),
		slog.String("blueprint_id", string(id)),
	)
	err := r.softDeleteBlueprintRow(ctx, r.db, string(id))
	if errors.Is(err, pgutils.ErrNoAffectedRows) {
		l.WarnContext(ctx, "blueprint not found", slog.String("error", err.Error()))
		return fmt.Errorf("%w: %s", ports.ErrBlueprintNotFound, id)
	}
	if err != nil {
		l.ErrorContext(ctx, "failed to delete blueprint", slog.String("error", err.Error()))
		return err
	}
	return nil
}
