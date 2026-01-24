package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zhikh23/pgutils"

	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

func (r *Repository) Blueprint(ctx context.Context, id value.BlueprintID) (*entity.Blueprint, error) {
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
		return nil, fmt.Errorf("%w: %s", ports.ErrBlueprintNotFound, string(id))
	}
	if err != nil {
		return nil, err
	}
	return blueprint, nil
}

func (r *Repository) Blueprints(ctx context.Context, uid value.UserID) ([]*entity.Blueprint, error) {
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
		return nil, err
	}
	return blueprints, nil
}

func (r *Repository) SearchBlueprints(ctx context.Context, uid value.UserID, name string) ([]*entity.Blueprint, error) {
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
		return nil, err
	}
	return blueprints, nil
}
