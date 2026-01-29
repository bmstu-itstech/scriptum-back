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

func (r *Repository) Job(ctx context.Context, id value.JobID) (dto.Job, error) {
	var rJ readJobRow
	var rIFs []jobFieldRow
	var rOFs []jobFieldRow
	var rIVs []jobValueRow
	var rOVs []jobValueRow

	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		var err error
		rJ, err = r.selectReadJobRow(ctx, tx, string(id))
		if err != nil {
			return err
		}
		rIFs, err = r.selectJobInputFieldRows(ctx, tx, string(id))
		if err != nil {
			return err
		}
		rOFs, err = r.selectJobOutputFieldRows(ctx, tx, string(id))
		if err != nil {
			return err
		}
		rIVs, err = r.selectJobInputValueRows(ctx, tx, string(id))
		if err != nil {
			return err
		}
		rOVs, err = r.selectJobOutputValueRows(ctx, tx, string(id))
		if err != nil {
			return err
		}
		return nil
	})
	if errors.Is(err, sql.ErrNoRows) {
		return dto.Job{}, fmt.Errorf("%w: %s", ports.ErrJobNotFound, string(id))
	}
	if err != nil {
		return dto.Job{}, err
	}

	return readJobRowToDTO(rJ, rIFs, rOFs, rIVs, rOVs), nil
}

func (r *Repository) UserJobs(ctx context.Context, uid value.UserID) ([]dto.Job, error) {
	var rJs []readJobRow
	var rIFs map[string][]jobFieldRow
	var rOFs map[string][]jobFieldRow
	var rIVs map[string][]jobValueRow
	var rOVs map[string][]jobValueRow

	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		var err error
		rJs, err = r.selectUserReadJobRows(ctx, tx, string(uid))
		if err != nil {
			return err
		}
		ids := idsFromJobs(rJs)
		rIFs, err = r.selectJobsInputFieldsRows(ctx, tx, ids)
		if err != nil {
			return err
		}
		rOFs, err = r.selectJobsOutputFieldsRows(ctx, tx, ids)
		if err != nil {
			return err
		}
		rIVs, err = r.selectJobsInputValuesRows(ctx, tx, ids)
		if err != nil {
			return err
		}
		rOVs, err = r.selectJobsOutputValuesRows(ctx, tx, ids)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	js := make([]dto.Job, len(rJs))
	for i, rJ := range rJs {
		js[i] = readJobRowToDTO(rJ, rIFs[rJ.ID], rOFs[rJ.ID], rIVs[rJ.ID], rOVs[rJ.ID])
	}

	return js, nil
}

func (r *Repository) UserJobsWithState(
	ctx context.Context,
	uid value.UserID,
	state value.JobState,
) ([]dto.Job, error) {
	var rJs []readJobRow
	var rIFs map[string][]jobFieldRow
	var rOFs map[string][]jobFieldRow
	var rIVs map[string][]jobValueRow
	var rOVs map[string][]jobValueRow

	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		var err error
		rJs, err = r.selectUserReadJobRowsWithState(ctx, tx, string(uid), state.String())
		if err != nil {
			return err
		}
		ids := idsFromJobs(rJs)
		rIFs, err = r.selectJobsInputFieldsRows(ctx, tx, ids)
		if err != nil {
			return err
		}
		rOFs, err = r.selectJobsOutputFieldsRows(ctx, tx, ids)
		if err != nil {
			return err
		}
		rIVs, err = r.selectJobsInputValuesRows(ctx, tx, ids)
		if err != nil {
			return err
		}
		rOVs, err = r.selectJobsOutputValuesRows(ctx, tx, ids)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	js := make([]dto.Job, len(rJs))
	for i, rJ := range rJs {
		js[i] = readJobRowToDTO(rJ, rIFs[rJ.ID], rOFs[rJ.ID], rIVs[rJ.ID], rOVs[rJ.ID])
	}

	return js, nil
}

func idsFromJobs(rJs []readJobRow) []string {
	res := make([]string, len(rJs))
	for i, rJ := range rJs {
		res[i] = rJ.ID
	}
	return res
}
