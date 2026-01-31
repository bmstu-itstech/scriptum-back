package postgres

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zhikh23/pgutils"
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

func (r *Repository) selectBlueprintWithUserRow(ctx context.Context, qc sqlx.QueryerContext, blueprintID string) (blueprintWithUserRow, error) {
	var row blueprintWithUserRow
	err := pgutils.Get(ctx, qc, &row, `
		SELECT
			b.id,
			b.archive_id,
			b.name,
			b."desc",
			b.vis,
			b.owner_id,
			u.name AS owner_name,
			b.created_at
		FROM blueprint.blueprints b
		LEFT JOIN users u
			ON u.id = b.owner_id
			AND u.deleted_at IS NULL
		WHERE 
			b.id = $1
			AND b.deleted_at IS NULL
		`,
		blueprintID,
	)
	if err != nil {
		return blueprintWithUserRow{}, fmt.Errorf("select blueprint with user row: %w", err)
	}
	return row, nil
}

func (r *Repository) selectPublicAndUserBlueprintWithUserRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	userID string,
) ([]blueprintWithUserRow, error) {
	var rows []blueprintWithUserRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			b.id,
			b.archive_id,
			b.name,
			b."desc",
			b.vis,
			b.owner_id,
			u.name AS owner_name,
			b.created_at
		FROM blueprint.blueprints b
		LEFT JOIN users u
			ON u.id = b.owner_id
			AND u.deleted_at IS NULL
		WHERE
			b.deleted_at IS NULL
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

func (r *Repository) selectPublicAndUserBlueprintWithUserByNameRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	userID string,
	name string,
) ([]blueprintWithUserRow, error) {
	var rows []blueprintWithUserRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			b.id,
			b.archive_id,
			b.name,
			b."desc",
			b.vis,
			b.owner_id,
			u.name AS owner_name,
			b.created_at
		FROM blueprint.blueprints b
		LEFT JOIN users u
			ON u.id = b.owner_id
			AND u.deleted_at IS NULL
		WHERE
			b.deleted_at IS NULL
			AND b.name ILIKE $2
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
		return nil, fmt.Errorf("select public and user blueprint with user rows: %w", err)
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

func (r *Repository) softDeleteBlueprintJobRows(ctx context.Context, ec sqlx.ExecerContext, blueprintID string) error {
	_, err := pgutils.Exec(ctx, ec, `
		UPDATE job.jobs
		SET
			deleted_at = NOW()
		WHERE id = $1
		`,
		blueprintID,
	)
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

func (r *Repository) selectBlueprintsInputFieldRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	blueprintIDs []string,
) (map[string][]blueprintFieldRow, error) {
	if len(blueprintIDs) == 0 {
		return map[string][]blueprintFieldRow{}, nil
	}
	query, args, err := sqlx.In(`
		SELECT
			blueprint_id, 
			index, 
			type, 
			name, 
			"desc", 
			unit
		FROM blueprint.input_fields
		WHERE
			blueprint_id IN (?)
		ORDER BY index
		`,
		blueprintIDs,
	)
	if err != nil {
		return nil, fmt.Errorf("sqlx.In: %w", err)
	}
	query = r.db.Rebind(query)

	var rows []blueprintFieldRow
	err = pgutils.Select(ctx, qc, &rows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("pgutils.Select: %w", err)
	}
	return mapBlueprintFieldRows(rows), nil
}

func mapBlueprintFieldRows(bs []blueprintFieldRow) map[string][]blueprintFieldRow {
	m := make(map[string][]blueprintFieldRow)
	for _, row := range bs {
		key := row.BlueprintID
		m[key] = append(m[key], row)
	}
	return m
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

func (r *Repository) selectBlueprintsOutputFieldRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	blueprintIDs []string,
) (map[string][]blueprintFieldRow, error) {
	if len(blueprintIDs) == 0 {
		return map[string][]blueprintFieldRow{}, nil
	}
	var rows []blueprintFieldRow
	query, args, err := sqlx.In(`
		SELECT
			blueprint_id, 
			index, 
			type, 
			name, 
			"desc", 
			unit
		FROM blueprint.output_fields
		WHERE
			blueprint_id IN (?)
		ORDER BY index
		`,
		blueprintIDs,
	)
	if err != nil {
		return nil, fmt.Errorf("sqlx.In: %w", err)
	}
	query = r.db.Rebind(query)
	err = pgutils.Select(ctx, qc, &rows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("pgutils.Select: %w", err)
	}
	return mapBlueprintFieldRows(rows), nil
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
		WHERE 
			id = $1
			AND deleted_at IS NULL
		`,
		jobID,
	)
	if err != nil {
		return jobRow{}, fmt.Errorf("select job row: %w", err)
	}
	return row, nil
}

func (r *Repository) selectReadJobRow(ctx context.Context, qc sqlx.QueryerContext, jobID string) (readJobRow, error) {
	var row readJobRow
	err := pgutils.Get(ctx, qc, &row, `
		SELECT
			j.id, 
			j.owner_id,
			j.blueprint_id, 
			b.name AS blueprint_name,
			j.state, 
			j.created_at, 
			j.started_at, 
			j.result_code, 
			j.result_msg, 
			j.finished_at
		FROM job.jobs j
		JOIN blueprint.blueprints b 
			ON j.blueprint_id = b.id
			AND b.deleted_at IS NULL
		WHERE 
			j.id = $1
			AND j.deleted_at IS NULL
		`,
		jobID,
	)
	return row, err
}

func (r *Repository) selectUserReadJobRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	uid string,
) ([]readJobRow, error) {
	var rows []readJobRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			j.id, 
			j.owner_id,
			j.blueprint_id, 
			b.name AS blueprint_name,
			j.state, 
			j.created_at, 
			j.started_at, 
			j.result_code, 
			j.result_msg, 
			j.finished_at
		FROM job.jobs j
		JOIN blueprint.blueprints b 
			ON j.blueprint_id = b.id
			AND b.deleted_at IS NULL
		WHERE 
			j.owner_id = $1
			AND j.deleted_at IS NULL
		ORDER BY created_at DESC
		`,
		uid,
	)
	if err != nil {
		return nil, fmt.Errorf("select user job rows: %w", err)
	}
	return rows, nil
}

func (r *Repository) selectUserReadJobRowsWithState(
	ctx context.Context,
	qc sqlx.QueryerContext,
	uid string,
	state string,
) ([]readJobRow, error) {
	var rows []readJobRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			j.id, 
			j.owner_id,
			j.blueprint_id, 
			b.name AS blueprint_name,
			j.state, 
			j.created_at, 
			j.started_at, 
			j.result_code, 
			j.result_msg, 
			j.finished_at
		FROM job.jobs j
		JOIN blueprint.blueprints b 
			ON j.blueprint_id = b.id
			AND b.deleted_at IS NULL
		WHERE 
			j.owner_id = $1
			AND state = $2
			AND j.deleted_at IS NULL
		ORDER BY created_at DESC
		`,
		uid,
		state,
	)
	if err != nil {
		return nil, fmt.Errorf("select user job rows: %w", err)
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
		WHERE 
			id = :id
			AND deleted_at IS NULL
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

func (r *Repository) selectJobsInputValuesRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	jobIDs []string,
) (map[string][]jobValueRow, error) {
	if len(jobIDs) == 0 {
		return map[string][]jobValueRow{}, nil
	}
	query, args, err := sqlx.In(`
		SELECT
			job_id, 
			index, 
			type, 
			value
		FROM job.input_values
		WHERE job_id IN (?)
		ORDER BY index
		`,
		jobIDs,
	)
	if err != nil {
		return nil, fmt.Errorf("sqlx.In: %w", err)
	}
	query = r.db.Rebind(query)

	var rows []jobValueRow
	err = pgutils.Select(ctx, qc, &rows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("pgutils.Select: %w", err)
	}

	return mapJobValueRows(rows), nil
}

func mapJobValueRows(rs []jobValueRow) map[string][]jobValueRow {
	m := make(map[string][]jobValueRow)
	for _, r := range rs {
		key := r.JobID
		m[key] = append(m[key], r)
	}
	return m
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

func (r *Repository) selectJobsOutputValuesRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	jobIDs []string,
) (map[string][]jobValueRow, error) {
	if len(jobIDs) == 0 {
		return map[string][]jobValueRow{}, nil
	}
	query, args, err := sqlx.In(`
		SELECT
			job_id, 
			index, 
			type, 
			value
		FROM job.output_values
		WHERE job_id IN (?)
		ORDER BY index
		`,
		jobIDs,
	)
	if err != nil {
		return nil, fmt.Errorf("sqlx.In: %w", err)
	}
	query = r.db.Rebind(query)

	var rows []jobValueRow
	err = pgutils.Select(ctx, qc, &rows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("pgutils.Select: %w", err)
	}

	return mapJobValueRows(rows), nil
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

func (r *Repository) selectJobInputFieldRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	jobID string,
) ([]jobFieldRow, error) {
	var rows []jobFieldRow
	err := pgutils.Select(ctx, qc, &rows, `
		SELECT
			j.id AS job_id,
			bif.index,
			bif.type,
			bif.name,
			bif."desc",
			bif.unit
		FROM blueprint.input_fields bif
		JOIN job.jobs j
			ON j.blueprint_id = bif.blueprint_id
		WHERE j.id = $1
		`,
		jobID,
	)
	if err != nil {
		return nil, fmt.Errorf("pgutils.Select: %w", err)
	}
	return rows, nil
}

func (r *Repository) selectJobsInputFieldsRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	jobIDs []string,
) (map[string][]jobFieldRow, error) {
	if len(jobIDs) == 0 {
		return map[string][]jobFieldRow{}, nil
	}
	query, args, err := sqlx.In(`
		SELECT
			j.id AS job_id, 
			bif.index,
			bif.type, 
			bif.name,
			bif."desc",
			bif.unit
		FROM blueprint.input_fields bif
		JOIN job.jobs j
			ON j.blueprint_id = bif.blueprint_id
		WHERE j.id IN (?)
		ORDER BY index
		`,
		jobIDs,
	)
	if err != nil {
		return nil, fmt.Errorf("sqlx.In: %w", err)
	}
	query = r.db.Rebind(query)

	var rows []jobFieldRow
	err = pgutils.Select(ctx, qc, &rows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("pgutils.Select: %w", err)
	}

	return mapJobFieldsRows(rows), nil
}

func (r *Repository) selectJobOutputFieldsRows(
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
		`,
		jobID,
	)
	if err != nil {
		return nil, fmt.Errorf("pgutils.Select: %w", err)
	}
	return rows, nil
}

func (r *Repository) selectJobsOutputFieldsRows(
	ctx context.Context,
	qc sqlx.QueryerContext,
	jobIDs []string,
) (map[string][]jobFieldRow, error) {
	if len(jobIDs) == 0 {
		return map[string][]jobFieldRow{}, nil
	}
	query, args, err := sqlx.In(`
		SELECT
			job_id, 
			index, 
			type, 
			name,
			"desc",
			unit
		FROM job.output_fields
		WHERE job_id IN (?)
		ORDER BY index
		`,
		jobIDs,
	)
	if err != nil {
		return nil, fmt.Errorf("sqlx.In: %w", err)
	}
	query = r.db.Rebind(query)

	var rows []jobFieldRow
	err = pgutils.Select(ctx, qc, &rows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("pgutils.Select: %w", err)
	}

	return mapJobFieldsRows(rows), nil
}

func mapJobFieldsRows(rows []jobFieldRow) map[string][]jobFieldRow {
	m := make(map[string][]jobFieldRow)
	for _, r := range rows {
		key := r.JobID
		m[key] = append(m[key], r)
	}
	return m
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

func (r *Repository) insertUserRow(
	ctx context.Context,
	ec sqlx.ExtContext,
	row userRow,
) error {
	err := pgutils.RequireAffected(pgutils.NamedExec(ctx, ec, `
		INSERT INTO users (
			id, 
			email, 
			name, 
			role, 
			passhash, 
			created_at
		)
		VALUES (
			:id,
			:email,
			:name,
			:role,
			:passhash,
			:created_at
		)
		`,
		row,
	))
	if err != nil {
		return fmt.Errorf("insert user row: %w", err)
	}
	return nil
}

func (r *Repository) updateUserRow(ctx context.Context, ec sqlx.ExtContext, row userRow) error {
	err := pgutils.RequireAffected(pgutils.NamedExec(ctx, ec, `
		UPDATE users
		SET
			email = :email,
			passhash = :passhash,
			name = :name,
			role = :role
		WHERE 
			id = :id
			AND deleted_at IS NULL
		`,
		row,
	))
	if err != nil {
		return fmt.Errorf("update user row: %w", err)
	}
	return nil
}

func (r *Repository) softDeleteUserRow(ctx context.Context, ec sqlx.ExecerContext, uid string) error {
	err := pgutils.RequireAffected(pgutils.Exec(ctx, ec, `
		UPDATE users
		SET
			deleted_at = NOW()
		WHERE id = $1
		`,
		uid,
	))
	if err != nil {
		return fmt.Errorf("failed to soft delete user row: %w", err)
	}
	return nil
}
