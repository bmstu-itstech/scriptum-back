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

func (r *Repository) Blueprint(ctx context.Context, id value.BlueprintID) (*entity.Blueprint, error) {
	l := r.l.With(
		slog.String("op", "postgres.Repository.Blueprint"),
		slog.String("blueprint_id", string(id)),
	)

	var blueprint *entity.Blueprint
	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		rB, err := r.selectBlueprintRow(ctx, tx, string(id))
		if err != nil {
			return err
		}
		rIn, err := r.selectBlueprintInputFieldRows(ctx, tx, string(id))
		if err != nil {
			return err
		}
		rOut, err := r.selectBlueprintOutputFieldRows(ctx, tx, string(id))
		if err != nil {
			return err
		}
		blueprint, err = blueprintRowToDomain(rB, rIn, rOut)
		if err != nil {
			return err
		}
		return nil
	})
	if errors.Is(err, sql.ErrNoRows) {
		l.WarnContext(ctx, "blueprint not found", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %s", ports.ErrBlueprintNotFound, string(id))
	}
	if err != nil {
		l.ErrorContext(ctx, "failed to execute transaction", slog.String("error", err.Error()))
		return nil, err
	}
	return blueprint, nil
}

func (r *Repository) Blueprints(ctx context.Context, uid value.UserID) ([]*entity.Blueprint, error) {
	l := r.l.With(
		slog.String("op", "postgres.Repository.Blueprints"),
		slog.String("user_id", string(uid)),
	)

	blueprints := make([]*entity.Blueprint, 0)
	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		rBs, err := r.selectPublicAndUserBlueprintRows(ctx, tx, string(uid))
		if err != nil {
			return err
		}
		for _, rB := range rBs {
			rIn, err2 := r.selectBlueprintInputFieldRows(ctx, tx, rB.ID)
			if err2 != nil {
				return err2
			}
			rOut, err2 := r.selectBlueprintOutputFieldRows(ctx, tx, rB.ID)
			if err2 != nil {
				return err2
			}
			blueprint, err2 := blueprintRowToDomain(rB, rIn, rOut)
			if err2 != nil {
				return err2
			}
			blueprints = append(blueprints, blueprint)
		}
		return nil
	})
	if err != nil {
		l.ErrorContext(ctx, "failed to execute transaction", slog.String("error", err.Error()))
		return nil, err
	}
	return blueprints, nil
}

func (r *Repository) SearchBlueprints(ctx context.Context, uid value.UserID, name string) ([]*entity.Blueprint, error) {
	l := r.l.With(
		slog.String("op", "postgres.Repository.SearchBlueprints"),
		slog.String("user_id", string(uid)),
		slog.String("name", name),
	)

	blueprints := make([]*entity.Blueprint, 0)
	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		rBs, err := r.selectPublicAndUserBlueprintByNameRows(ctx, tx, string(uid), name)
		if err != nil {
			return err
		}
		for _, rB := range rBs {
			rIn, err2 := r.selectBlueprintInputFieldRows(ctx, tx, rB.ID)
			if err2 != nil {
				return err2
			}
			rOut, err2 := r.selectBlueprintOutputFieldRows(ctx, tx, rB.ID)
			if err2 != nil {
				return err2
			}
			blueprint, err2 := blueprintRowToDomain(rB, rIn, rOut)
			if err2 != nil {
				return err2
			}
			blueprints = append(blueprints, blueprint)
		}
		return nil
	})
	if err != nil {
		l.ErrorContext(ctx, "failed to execute transaction", slog.String("error", err.Error()))
		return nil, err
	}
	return blueprints, nil
}
