package postgres

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/zhikh23/pgutils"

	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

func (r *Repository) SaveBox(ctx context.Context, box *entity.Box) error {
	l := r.l.With(
		slog.String("op", "postgres.Repository.SaveBox"),
		slog.String("box_id", string(box.ID())),
	)
	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		rBox := boxRowFromDomain(box)
		if err := r.insertBoxRow(ctx, tx, rBox); err != nil {
			return err
		}
		rIn := boxFieldRowsFromDomain(box.In(), box.ID())
		if err := r.insertBoxInputFieldRows(ctx, tx, rIn); err != nil {
			return err
		}
		rOut := boxFieldRowsFromDomain(box.Out(), box.ID())
		if err := r.insertBoxOutputFieldRows(ctx, tx, rOut); err != nil {
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

func (r *Repository) DeleteBox(ctx context.Context, id value.BoxID) error {
	l := r.l.With(
		slog.String("op", "postgres.Repository.DeleteBox"),
		slog.String("box_id", string(id)),
	)
	err := r.softDeleteBoxRow(ctx, r.db, string(id))
	if errors.Is(err, pgutils.ErrNoAffectedRows) {
		l.WarnContext(ctx, "box not found", slog.String("error", err.Error()))
		return ports.ErrBoxNotFound
	}
	if err != nil {
		l.ErrorContext(ctx, "failed to delete box", slog.String("error", err.Error()))
		return err
	}
	return nil
}
