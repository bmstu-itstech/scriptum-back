package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zhikh23/pgutils"

	"github.com/bmstu-itstech/scriptum-back/internal/app/dto"
	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

func (r *Repository) BlueprintWithUser(ctx context.Context, id value.BlueprintID) (dto.BlueprintWithUser, error) {
	var blueprint dto.BlueprintWithUser
	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		rB, err := r.selectBlueprintWithUserRow(ctx, tx, string(id))
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
		blueprint = blueprintWithUserRowToDTO(rB, rIn, rOut)
		return nil
	})
	if errors.Is(err, sql.ErrNoRows) {
		return dto.BlueprintWithUser{}, fmt.Errorf("%w: %s", ports.ErrBlueprintNotFound, string(id))
	}
	if err != nil {
		return dto.BlueprintWithUser{}, err
	}
	return blueprint, nil
}

func (r *Repository) BlueprintsWithUsers(ctx context.Context, uid value.UserID) ([]dto.BlueprintWithUser, error) {
	blueprints := make([]dto.BlueprintWithUser, 0)
	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		rBs, err := r.selectPublicAndUserBlueprintWithUserRows(ctx, tx, string(uid))
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
			blueprint := blueprintWithUserRowToDTO(rB, rIn, rOut)
			blueprints = append(blueprints, blueprint)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return blueprints, nil
}

func (r *Repository) SearchBlueprintsWithUsers(ctx context.Context, uid value.UserID, name string) ([]dto.BlueprintWithUser, error) {
	blueprints := make([]dto.BlueprintWithUser, 0)
	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		rBs, err := r.selectPublicAndUserBlueprintWithUserByNameRows(ctx, tx, string(uid), name)
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
			blueprint := blueprintWithUserRowToDTO(rB, rIn, rOut)
			blueprints = append(blueprints, blueprint)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return blueprints, nil
}
