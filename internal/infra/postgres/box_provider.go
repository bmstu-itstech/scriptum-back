package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/zhikh23/pgutils"

	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

func (r *Repository) Box(ctx context.Context, id value.BoxID) (*entity.Box, error) {
	l := r.l.With(
		slog.String("op", "postgres.Repository.Box"),
		slog.String("box_id", string(id)),
	)

	var box *entity.Box
	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		rBox, err := r.selectBoxRow(ctx, tx, string(id))
		if err != nil {
			return err
		}
		rIn, err := r.selectBoxInputFieldRows(ctx, tx, string(id))
		if err != nil {
			return err
		}
		rOut, err := r.selectBoxOutputFieldRows(ctx, tx, string(id))
		if err != nil {
			return err
		}
		box, err = boxRowToDomain(rBox, rIn, rOut)
		if err != nil {
			return err
		}
		return nil
	})
	if errors.Is(err, sql.ErrNoRows) {
		l.WarnContext(ctx, "box not found", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %s", ports.ErrBoxNotFound, string(id))
	}
	if err != nil {
		l.ErrorContext(ctx, "failed to execute transaction", slog.String("error", err.Error()))
		return nil, err
	}
	return box, nil
}

func (r *Repository) Boxes(ctx context.Context, uid value.UserID) ([]*entity.Box, error) {
	l := r.l.With(
		slog.String("op", "postgres.Repository.Boxes"),
		slog.Int64("user_id", int64(uid)),
	)

	boxes := make([]*entity.Box, 0)
	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		rBoxes, err := r.selectPublicAndUserBoxRows(ctx, tx, int64(uid))
		if err != nil {
			return err
		}
		for _, rBox := range rBoxes {
			rIn, err2 := r.selectBoxInputFieldRows(ctx, tx, rBox.ID)
			if err2 != nil {
				return err2
			}
			rOut, err2 := r.selectBoxOutputFieldRows(ctx, tx, rBox.ID)
			if err2 != nil {
				return err2
			}
			box, err2 := boxRowToDomain(rBox, rIn, rOut)
			if err2 != nil {
				return err2
			}
			boxes = append(boxes, box)
		}
		return nil
	})
	if err != nil {
		l.ErrorContext(ctx, "failed to execute transaction", slog.String("error", err.Error()))
		return nil, err
	}
	return boxes, nil
}

func (r *Repository) SearchBoxes(ctx context.Context, uid value.UserID, name string) ([]*entity.Box, error) {
	l := r.l.With(
		slog.String("op", "postgres.Repository.SearchBoxes"),
		slog.Int64("user_id", int64(uid)),
		slog.String("name", name),
	)

	boxes := make([]*entity.Box, 0)
	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		rBoxes, err := r.selectPublicAndUserBoxByNameRows(ctx, tx, int64(uid), name)
		if err != nil {
			return err
		}
		for _, rBox := range rBoxes {
			rIn, err2 := r.selectBoxInputFieldRows(ctx, tx, rBox.ID)
			if err2 != nil {
				return err2
			}
			rOut, err2 := r.selectBoxOutputFieldRows(ctx, tx, rBox.ID)
			if err2 != nil {
				return err2
			}
			box, err2 := boxRowToDomain(rBox, rIn, rOut)
			if err2 != nil {
				return err2
			}
			boxes = append(boxes, box)
		}
		return nil
	})
	if err != nil {
		l.ErrorContext(ctx, "failed to execute transaction", slog.String("error", err.Error()))
		return nil, err
	}
	return boxes, nil
}
