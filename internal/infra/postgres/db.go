package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zhikh23/pgutils"

	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
)

func (r *Repository) selectBoxRow(ctx context.Context, qc sqlx.QueryerContext, boxID string) (boxRow, error) {
	var row boxRow
	err := pgutils.Get(ctx, qc, &row, `
		SELECT
			id,
			owner_id,
			archive_id,
			name,
			"desc",
			vis,
			created_at
		FROM box.boxes
		WHERE 
			id = $1
			AND deleted_at IS NULL
		`,
		boxID,
	)
	if err != nil {
		return boxRow{}, fmt.Errorf("select box row: %w", err)
	}
	return row, nil
}

func (r *Repository) selectPublicAndUserBoxRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	userID int64,
) ([]boxRow, error) {
	var rows []boxRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			id,
			owner_id,
			archive_id,
			name,
			"desc",
			vis,
			created_at
		FROM box.boxes
		WHERE
			deleted_at IS NULL
			AND (
			    vis = 'public'
			    OR owner_id = $1
			)
		ORDER BY created_at DESC
		`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("select public and user box rows: %w", err)
	}
	return rows, nil
}

func (r *Repository) selectPublicAndUserBoxByNameRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	userID int64,
	name string,
) ([]boxRow, error) {
	var rows []boxRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			id,
			owner_id,
			archive_id,
			name,
			"desc",
			vis,
			created_at
		FROM box.boxes
		WHERE
			deleted_at IS NULL
			AND name ILIKE $2
			AND (
			    vis = 'public'
			    OR owner_id = $1
			)
		ORDER BY created_at DESC
		`,
		userID,
		name,
	)
	if err != nil {
		return nil, fmt.Errorf("select public and user box rows: %w", err)
	}
	return rows, nil
}

func (r *Repository) insertBoxRow(ctx context.Context, ec sqlx.ExecerContext, row boxRow) error {
	err := pgutils.RequireAffected(pgutils.Exec(ctx, ec, `
		INSERT INTO box.boxes (
			id,
			owner_id,
			archive_id,
			name,
			"desc",
			vis,
			created_at
		)
		VALUES (
			:id, 
			:owner_id, 
			:archive_id, 
			:name, 
			:desc, 
			:vis, 
			:created_at
		)
		`,
		row,
	))
	if pgutils.IsUniqueViolationError(err) {
		return ports.ErrJobAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("upsert box row: %w", err)
	}
	return nil
}

func (r *Repository) softDeleteBoxRow(ctx context.Context, ec sqlx.ExecerContext, boxID string) error {
	err := pgutils.RequireAffected(pgutils.Exec(ctx, ec, `
		UPDATE box.boxes
		SET
			deleted_at = NOW()
		WHERE id = $1
		`,
		boxID,
	))
	if err != nil {
		return fmt.Errorf("failed to soft delete box row: %w", err)
	}
	return nil
}

func (r *Repository) selectBoxInputFieldRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	boxID string,
) ([]boxFieldRow, error) {
	var rows []boxFieldRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			box_id, 
			index, 
			type, 
			name, 
			"desc", 
			unit
		FROM box.input_fields
		WHERE box_id = $1
		ORDER BY index
		`,
		boxID,
	)
	if err != nil {
		return nil, fmt.Errorf("select box input fields rows: %w", err)
	}
	return rows, nil
}

func (r *Repository) insertBoxInputFieldRows(
	ctx context.Context,
	ec sqlx.ExecerContext,
	rows []boxFieldRow,
) error {
	err := pgutils.RequireAffected(pgutils.Exec(ctx, ec, `
		INSERT INTO box.input_fields (
			box_id, 
			index, 
			type, 
			name, 
			"desc", 
			unit
		)
		VALUES (
		    :box_id,
			:index,
			:type,
			:name,
			:desc,
			:unit
		)	
		`,
		rows,
	))
	if err != nil {
		return fmt.Errorf("insert box input field rows: %w", err)
	}
	return nil
}

func (r *Repository) selectBoxOutputFieldRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	boxID string,
) ([]boxFieldRow, error) {
	var rows []boxFieldRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			box_id, 
			index, 
			type, 
			name, 
			"desc", 
			unit
		FROM box.output_fields
		WHERE box_id = $1
		ORDER BY index
		`,
		boxID,
	)
	if err != nil {
		return nil, fmt.Errorf("select box output fields rows: %w", err)
	}
	return rows, nil
}

func (r *Repository) insertBoxOutputFieldRows(
	ctx context.Context,
	ec sqlx.ExecerContext,
	rows []boxFieldRow,
) error {
	err := pgutils.RequireAffected(pgutils.Exec(ctx, ec, `
		INSERT INTO box.output_fields (
			box_id, 
			index, 
			type, 
			name, 
			"desc", 
			unit
		)
		VALUES (
		    :box_id,
			:index,
			:type,
			:name,
			:desc,
			:unit
		)	
		`,
		rows,
	))
	if err != nil {
		return fmt.Errorf("insert box output field rows: %w", err)
	}
	return nil
}

func (r *Repository) selectJobRow(ctx context.Context, qc sqlx.QueryerContext, jobID string) (jobRow, error) {
	var row jobRow
	err := pgutils.Select(ctx, qc, &row, `
		SELECT
			id, 
			box_id, 
			archive_id, 
			owner_id, 
			state, 
			created_at, 
			started_at, 
			result_code, 
			result_msg, 
			finished_at
		FROM job.jobs
		WHERE id = $1
		`,
		jobID,
	)
	if err != nil {
		return jobRow{}, fmt.Errorf("select job row: %w", err)
	}
	return row, nil
}

func (r *Repository) selectUserJobRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	uid int64,
) ([]jobRow, error) {
	var rows []jobRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			id, 
			box_id, 
			archive_id, 
			owner_id, 
			state, 
			created_at, 
			started_at, 
			result_code, 
			result_msg, 
			finished_at
		FROM job.jobs
		WHERE owner_id = $1
		ORDER BY created_at DESC
		`,
		uid,
	)
	if err != nil {
		return nil, fmt.Errorf("select user job rows: %w", err)
	}
	return rows, nil
}

func (r *Repository) selectUserJobRowsWithState(
	ctx context.Context,
	qc sqlx.QueryerContext,
	uid int64,
	state string,
) ([]jobRow, error) {
	var rows []jobRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			id, 
			box_id, 
			archive_id, 
			owner_id, 
			state, 
			created_at, 
			started_at, 
			result_code, 
			result_msg, 
			finished_at
		FROM job.jobs
		WHERE 
			owner_id = $1
			AND state = $2
		ORDER BY created_at DESC
		`,
		uid,
		state,
	)
	if err != nil {
		return nil, fmt.Errorf("select user job rows with state: %w", err)
	}
	return rows, nil
}

func (r *Repository) insertJobRow(ctx context.Context, ec sqlx.ExecerContext, row jobRow) error {
	err := pgutils.RequireAffected(pgutils.Exec(ctx, ec, `
		INSERT INTO job.jobs (
		    id, 
			box_id, 
			archive_id, 
			owner_id, 
			state, 
			created_at, 
			started_at, 
			result_code, 
			result_msg, 
			finished_at
		) 
		VALUES (
			:id,
			:box_id,
			:archive_id,
			:owner_id,
			:state,
			:created_at,
			:started_at,
			:result_code,
			:result_msg,
			:finished_at
		)
		`,
		row,
	))
	if pgutils.IsUniqueViolationError(err) {
		return ports.ErrJobAlreadyExists
	}
	if err != nil {
		return fmt.Errorf("insert job row: %w", err)
	}
	return nil
}

func (r *Repository) selectJobInputValueRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	jobID string,
) ([]jobValueRow, error) {
	var rows []jobValueRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			job_id, 
			index, 
			type, 
			value
		FROM job.input_values
		WHERE job_id = $1
		ORDER BY index
		`,
		jobID,
	)
	if err != nil {
		return nil, fmt.Errorf("select job input value rows: %w", err)
	}
	return rows, nil
}

func (r *Repository) insertJobInputValueRows(
	ctx context.Context,
	ec sqlx.ExecerContext,
	rows []jobValueRow,
) error {
	err := pgutils.RequireAffected(pgutils.Exec(ctx, ec, `
		INSERT INTO job.input_values (
		    job_id, 
			index, 
			type, 
			value
		) 
		VALUES (
			:job_id, 
			:index,
			:type,
			:value
		)
		`,
		rows,
	))
	if err != nil {
		return fmt.Errorf("insert job input value rows: %w", err)
	}
	return nil
}

func (r *Repository) selectJobOutputValueRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	jobID string,
) ([]jobValueRow, error) {
	var rows []jobValueRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			job_id, 
			index, 
			type, 
			value
		FROM job.output_values
		WHERE job_id = $1
		ORDER BY index
		`,
		jobID,
	)
	if err != nil {
		return nil, fmt.Errorf("select job output value rows: %w", err)
	}
	return rows, nil
}

func (r *Repository) insertJobOutputValueRows(
	ctx context.Context,
	ec sqlx.ExecerContext,
	rows []jobValueRow,
) error {
	err := pgutils.RequireAffected(pgutils.Exec(ctx, ec, `
		INSERT INTO job.output_values (
		    job_id, 
			index, 
			type, 
			value
		) 
		VALUES (
			:job_id, 
			:index,
			:type,
			:value
		)
		`,
		rows,
	))
	if err != nil {
		return fmt.Errorf("insert job output value rows: %w", err)
	}
	return nil
}

func (r *Repository) selectJobOutputFieldRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	jobID string,
) ([]jobFieldRow, error) {
	var rows []jobFieldRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			job_id, 
			index, 
			type, 
			name, 
			"desc", 
			unit
		FROM job.output_fields
		WHERE job_id = $1
		ORDER BY index
		`,
		jobID,
	)
	if err != nil {
		return nil, fmt.Errorf("select job output fields rows: %w", err)
	}
	return rows, nil
}

func (r *Repository) insertJobOutputFieldRows(
	ctx context.Context,
	ec sqlx.ExecerContext,
	rows []jobFieldRow,
) error {
	err := pgutils.RequireAffected(pgutils.Exec(ctx, ec, `
		INSERT INTO job.output_fields (
			job_id, 
			index, 
			type, 
			name, 
			"desc", 
			unit
		)
		VALUES (
		    :job_id,
			:index,
			:type,
			:name,
			:desc,
			:unit
		)	
		`,
		rows,
	))
	if err != nil {
		return fmt.Errorf("insert job output field rows: %w", err)
	}
	return nil
}
