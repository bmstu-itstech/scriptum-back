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
	var rB blueprintWithUserRow
	var rIs []blueprintFieldRow
	var rOs []blueprintFieldRow

	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		var err error
		rB, err = r.selectBlueprintWithUserRow(ctx, tx, string(id))
		if err != nil {
			return err
		}
		rIs, err = r.selectBlueprintInputFieldRows(ctx, tx, string(id))
		if err != nil {
			return err
		}
		rOs, err = r.selectBlueprintOutputFieldRows(ctx, tx, string(id))
		if err != nil {
			return err
		}
		return nil
	})
	if errors.Is(err, sql.ErrNoRows) {
		return dto.BlueprintWithUser{}, fmt.Errorf("%w: %s", ports.ErrBlueprintNotFound, string(id))
	}
	if err != nil {
		return dto.BlueprintWithUser{}, err
	}

	return blueprintWithUserRowToDTO(rB, rIs, rOs), nil
}

func (r *Repository) BlueprintsWithUsers(ctx context.Context, uid value.UserID) ([]dto.BlueprintWithUser, error) {
	var rBs []blueprintWithUserRow
	var rIs map[string][]blueprintFieldRow
	var rOs map[string][]blueprintFieldRow

	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		var err error
		rBs, err = r.selectPublicAndUserBlueprintWithUserRows(ctx, tx, string(uid))
		if err != nil {
			return err
		}
		ids := idsFromBlueprints(rBs)
		rIs, err = r.selectBlueprintsInputFieldRows(ctx, tx, ids)
		if err != nil {
			return err
		}
		rOs, err = r.selectBlueprintsOutputFieldRows(ctx, tx, ids)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	bs := make([]dto.BlueprintWithUser, len(rBs))
	for i, rB := range rBs {
		bs[i] = blueprintWithUserRowToDTO(rB, rIs[rB.ID], rOs[rB.ID])
	}

	return bs, nil
}

func (r *Repository) SearchBlueprintsWithUsers(ctx context.Context, uid value.UserID, name string) ([]dto.BlueprintWithUser, error) {
	var rBs []blueprintWithUserRow
	var rIs map[string][]blueprintFieldRow
	var rOs map[string][]blueprintFieldRow

	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		var err error
		rBs, err = r.selectPublicAndUserBlueprintWithUserByNameRows(ctx, tx, string(uid), name)
		if err != nil {
			return err
		}
		ids := idsFromBlueprints(rBs)
		rIs, err = r.selectBlueprintsInputFieldRows(ctx, tx, ids)
		if err != nil {
			return err
		}
		rOs, err = r.selectBlueprintsOutputFieldRows(ctx, tx, ids)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	bs := make([]dto.BlueprintWithUser, len(rBs))
	for i, rB := range rBs {
		bs[i] = blueprintWithUserRowToDTO(rB, rIs[rB.ID], rOs[rB.ID])
	}

	return bs, nil
}

func idsFromBlueprints(bRs []blueprintWithUserRow) []string {
	res := make([]string, len(bRs))
	for i, bR := range bRs {
		res[i] = bR.ID
	}
	return res
}
