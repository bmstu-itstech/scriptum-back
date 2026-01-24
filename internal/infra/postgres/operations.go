package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zhikh23/pgutils"

	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
)

func (r *Repository) selectBlueprintRow(ctx context.Context, qc sqlx.QueryerContext, blueprintID string) (blueprintRow, error) {
	var row blueprintRow
	err := pgutils.Get(ctx, qc, &row, `
		SELECT
			id,
			owner_id,
			archive_id,
			name,
			"desc",
			vis,
			created_at
		FROM blueprint.blueprints
		WHERE 
			id = $1
			AND deleted_at IS NULL
		`,
		blueprintID,
	)
	if err != nil {
		return blueprintRow{}, fmt.Errorf("select blueprint row: %w", err)
	}
	return row, nil
}

func (r *Repository) selectPublicAndUserBlueprintRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	userID string,
) ([]blueprintRow, error) {
	var rows []blueprintRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			id,
			owner_id,
			archive_id,
			name,
			"desc",
			vis,
			created_at
		FROM blueprint.blueprints
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
		return nil, fmt.Errorf("select public and user blueprint rows: %w", err)
	}
	return rows, nil
}

func (r *Repository) selectPublicAndUserBlueprintByNameRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	userID string,
	name string,
) ([]blueprintRow, error) {
	var rows []blueprintRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			id,
			owner_id,
			archive_id,
			name,
			"desc",
			vis,
			created_at
		FROM blueprint.blueprints
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
		"%"+name+"%",
	)
	if err != nil {
		return nil, fmt.Errorf("select public and user blueprint rows: %w", err)
	}
	return rows, nil
}

func (r *Repository) insertBlueprintRow(ctx context.Context, ec sqlx.ExtContext, row blueprintRow) error {
	err := pgutils.RequireAffected(pgutils.NamedExec(ctx, ec, `
		INSERT INTO blueprint.blueprints (
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
		return fmt.Errorf("upsert blueprint row: %w", err)
	}
	return nil
}

func (r *Repository) softDeleteBlueprintRow(ctx context.Context, ec sqlx.ExecerContext, blueprintID string) error {
	err := pgutils.RequireAffected(pgutils.Exec(ctx, ec, `
		UPDATE blueprint.blueprints
		SET
			deleted_at = NOW()
		WHERE id = $1
		`,
		blueprintID,
	))
	if err != nil {
		return fmt.Errorf("failed to soft delete blueprint row: %w", err)
	}
	return nil
}

func (r *Repository) selectBlueprintInputFieldRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	blueprintID string,
) ([]blueprintFieldRow, error) {
	var rows []blueprintFieldRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			blueprint_id, 
			index, 
			type, 
			name, 
			"desc", 
			unit
		FROM blueprint.input_fields
		WHERE blueprint_id = $1
		ORDER BY index
		`,
		blueprintID,
	)
	if err != nil {
		return nil, fmt.Errorf("select blueprint input fields rows: %w", err)
	}
	return rows, nil
}

func (r *Repository) insertBlueprintInputFieldRows(
	ctx context.Context,
	ec sqlx.ExtContext,
	rows []blueprintFieldRow,
) error {
	err := pgutils.RequireAffected(pgutils.NamedExec(ctx, ec, `
		INSERT INTO blueprint.input_fields (
			blueprint_id, 
			index, 
			type, 
			name, 
			"desc", 
			unit
		)
		VALUES (
		    :blueprint_id,
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
		return fmt.Errorf("insert blueprint input field rows: %w", err)
	}
	return nil
}

func (r *Repository) selectBlueprintOutputFieldRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	blueprintID string,
) ([]blueprintFieldRow, error) {
	var rows []blueprintFieldRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			blueprint_id, 
			index, 
			type, 
			name, 
			"desc", 
			unit
		FROM blueprint.output_fields
		WHERE blueprint_id = $1
		ORDER BY index
		`,
		blueprintID,
	)
	if err != nil {
		return nil, fmt.Errorf("select blueprint output fields rows: %w", err)
	}
	return rows, nil
}

func (r *Repository) insertBlueprintOutputFieldRows(
	ctx context.Context,
	ec sqlx.ExtContext,
	rows []blueprintFieldRow,
) error {
	err := pgutils.RequireAffected(pgutils.NamedExec(ctx, ec, `
		INSERT INTO blueprint.output_fields (
			blueprint_id, 
			index, 
			type, 
			name, 
			"desc", 
			unit
		)
		VALUES (
		    :blueprint_id,
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
		return fmt.Errorf("insert blueprint output field rows: %w", err)
	}
	return nil
}

func (r *Repository) selectJobRow(ctx context.Context, qc sqlx.QueryerContext, jobID string) (jobRow, error) {
	var row jobRow
	err := pgutils.Get(ctx, qc, &row, `
		SELECT
			id, 
			blueprint_id, 
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
	uid string,
) ([]jobRow, error) {
	var rows []jobRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			id, 
			blueprint_id, 
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
	uid string,
	state string,
) ([]jobRow, error) {
	var rows []jobRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			id, 
			blueprint_id, 
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

func (r *Repository) insertJobRow(ctx context.Context, ec sqlx.ExtContext, row jobRow) error {
	err := pgutils.RequireAffected(pgutils.NamedExec(ctx, ec, `
		INSERT INTO job.jobs (
		    id, 
			blueprint_id, 
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
			:blueprint_id,
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

func (r *Repository) updateJobRow(ctx context.Context, ec sqlx.ExtContext, row jobRow) error {
	err := pgutils.RequireAffected(pgutils.NamedExec(ctx, ec, `
		UPDATE job.jobs
		SET
			state = :state,
			started_at = :started_at,
			result_code = :result_code,
			result_msg = :result_msg,
			finished_at = :finished_at
		WHERE id = :id
		`,
		row,
	))
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
	ec sqlx.ExtContext,
	rows []jobValueRow,
) error {
	_, err := pgutils.NamedExec(ctx, ec, `
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
		ON CONFLICT (job_id, index)
		DO NOTHING
		`,
		rows,
	)
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
	ec sqlx.ExtContext,
	rows []jobValueRow,
) error {
	_, err := pgutils.NamedExec(ctx, ec, `
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
		ON CONFLICT (job_id, index)
		DO NOTHING
		`,
		rows,
	)
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
	ec sqlx.ExtContext,
	rows []jobFieldRow,
) error {
	_, err := pgutils.NamedExec(ctx, ec, `
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
		ON CONFLICT (job_id, index)
		DO NOTHING
		`,
		rows,
	)
	if err != nil {
		return fmt.Errorf("insert job output field rows: %w", err)
	}
	return nil
}

func (r *Repository) selectUserRow(
	ctx context.Context,
	qc sqlx.QueryerContext,
	userID string,
) (userRow, error) {
	var row userRow
	err := pgutils.Get(ctx, qc, &row, `
		SELECT
			id,
			email,
			name,
			passhash,
			role,
			created_at
		FROM public.users
		WHERE 
			id = $1
			AND deleted_at IS NULL
		`,
		userID,
	)
	if err != nil {
		return userRow{}, fmt.Errorf("select user row: %w", err)
	}
	return row, nil
}

func (r *Repository) selectUserRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
) ([]userRow, error) {
	var rows []userRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			id,
			email,
			name,
			passhash,
			role,
			created_at
		FROM public.users
		WHERE 
			deleted_at IS NULL
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("select user rows: %w", err)
	}
	return rows, nil
}

func (r *Repository) selectUserRowByEmail(
	ctx context.Context,
	qc sqlx.QueryerContext,
	email string,
) (userRow, error) {
	var row userRow
	err := pgutils.Get(ctx, qc, &row, `
		SELECT
			id,
			email,
			name,
			passhash,
			role,
			created_at
		FROM public.users
		WHERE 
			email = $1
			AND deleted_at IS NULL
		`,
		email,
	)
	if err != nil {
		return userRow{}, fmt.Errorf("select user row: %w", err)
	}
	return row, nil
}
