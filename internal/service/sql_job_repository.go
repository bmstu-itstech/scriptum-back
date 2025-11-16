package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/jmoiron/sqlx"
)

type JobRepo struct {
	db *sqlx.DB
	l  *slog.Logger
}

func NewJobRepository(db *sqlx.DB, l *slog.Logger) *JobRepo {
	return &JobRepo{
		db: db,
		l:  l,
	}
}

const createJobQuery = `
	INSERT INTO jobs (user_id, script_id)
	VALUES (:user_id, :script_id)
	RETURNING job_id
`

func (r *JobRepo) Create(ctx context.Context, job *scripts.JobPrototype) (*scripts.Job, error) {
	r.l.Info("create job", "job", *job)
	r.l.Debug("begining transaction")
	tx, err := r.db.BeginTxx(ctx, nil)
	r.l.Debug("transaction started", "err", err)
	if err != nil {
		r.l.Error("failed to start transaction", "err", err.Error())
		return nil, err
	}
	defer func() {
		r.l.Debug("transaction finished", "err", err)
		if err != nil {
			r.l.Error("failed to commit transaction", "err", err.Error())
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
			r.l.Debug("transaction committed", "err", err)
		}
	}()

	var jobID int64

	r.l.Debug("creating job", "job", *job)
	named, args, err := sqlx.Named(createJobQuery, convertJobPrototypeToDB(job))
	r.l.Debug("named query", "named", named, "args", args, "err", err)
	if err != nil {
		return nil, err
	}
	query := tx.Rebind(named)

	r.l.Debug("executing query", "query", query, "args", args)
	err = tx.QueryRowContext(ctx, query, args...).Scan(&jobID)
	r.l.Debug("query executed", "err", err)
	if err != nil {
		r.l.Error("failed to execute query", "err", err.Error())
		return nil, err
	}

	r.l.Debug("inserting values", "jobID", jobID, "scriptID", job.ScriptID(), "input", job.Input())
	if err := insertValuesTx(ctx, tx, jobID, int64(job.ScriptID()), job.Input(), "in"); err != nil {
		r.l.Error("failed to insert values", "err", err.Error())
		return nil, err
	}

	r.l.Debug("job created", "jobID", jobID)
	return job.Build(scripts.JobID(jobID))
}

const deleteJobQuery = `DELETE FROM jobs WHERE job_id = $1`

func (r *JobRepo) Delete(ctx context.Context, jobID scripts.JobID) error {
	r.l.Info("delete job", "jobID", jobID)
	r.l.Debug("deleting job", "ctx", ctx)
	result, err := r.db.ExecContext(ctx, deleteJobQuery, jobID)
	r.l.Debug("deleted job", "result", result, "err", err)
	if err != nil {
		r.l.Error("failed to delete job", "err", err.Error())
		return err
	}

	r.l.Debug("rows affected")
	rowsAffected, err := result.RowsAffected()
	r.l.Debug("rows affected", "rowsAffected", rowsAffected, "err", err)
	if err != nil {
		return err
	}

	r.l.Debug("check if no rows affected", "is", rowsAffected == 0)
	if rowsAffected == 0 {
		r.l.Error("failed to delete job", "jobID", jobID)
		return fmt.Errorf("%w Delete: cannot delete job with id: %d", scripts.ErrJobNotFound, jobID)
	}

	return nil
}

const getJobQuery = "SELECT * FROM jobs WHERE job_id=$1"

const paramsJobQuery = `
		SELECT p.field_id, f.field_type, p.value, f.param
		FROM job_params jp
		JOIN parameters p ON jp.parameter_id = p.parameter_id
		JOIN fields f ON p.field_id = f.field_id
		WHERE jp.job_id = $1
		ORDER BY p.parameter_id
	`

func (r *JobRepo) Job(ctx context.Context, jobID scripts.JobID) (*scripts.Job, error) {
	r.l.Info("get job", "jobID", jobID)
	var jobRow JobRow
	r.l.Debug("getting job", "ctx", ctx)
	err := r.db.GetContext(ctx, &jobRow, getJobQuery, jobID)
	r.l.Debug("got job", "jobRow", jobRow, "err", err)
	if err != nil {
		r.l.Error("failed to get job", "err", err.Error())
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w Job: cannot extract job with id: %d", scripts.ErrJobNotFound, jobID)
		}
		return nil, err
	}
	var paramRows []ValueRow
	r.l.Debug("getting job params", "ctx", ctx)
	if err := r.db.SelectContext(ctx, &paramRows, paramsJobQuery, jobID); err != nil {
		r.l.Error("failed to get job params", "err", err.Error())
		return nil, err
	}
	r.l.Debug("got job params", "paramRows", paramRows)

	r.l.Debug("getting job values", "ctx", ctx)
	inputValues, outputValues, err := getJobValues(ctx, r.db, jobRow.ID)
	r.l.Debug("got job values", "inputValues", inputValues, "outputValues", outputValues, "err", err)
	if err != nil {
		r.l.Error("failed to get job values", "err", err.Error())
		return nil, err
	}

	r.l.Debug("restoring job", "jobRow", jobRow)
	var result *scripts.Result
	r.l.Debug("needed to restore", "jobRow.StatusCode != nil", jobRow.StatusCode != nil, "jobRow.ClosedAt != nil", jobRow.ClosedAt != nil)
	if jobRow.StatusCode != nil && jobRow.ClosedAt != nil {
		r.l.Debug("restoring result")
		result = scripts.RestoreResult(outputValues, scripts.StatusCode(*jobRow.StatusCode), jobRow.ErrorMessage)
	}

	r.l.Debug("getting output fields", "ctx", ctx)
	var outFields []fieldRow
	err = r.db.SelectContext(ctx, &outFields, getFieldsQuery, jobRow.ScriptID, "out")
	r.l.Debug("got output fields", "outFields", outFields, "err", err)
	if err != nil {
		r.l.Error("failed to get output fields", "err", err.Error())
		return nil, err
	}
	r.l.Debug("converting output fields", "outFields", outFields)
	outputs, err := convertFieldRowsToDomain(outFields)
	r.l.Debug("converted output fields", "outputs", outputs, "err", err)
	if err != nil {
		r.l.Error("failed to convert output fields", "err", err.Error())
		return nil, err
	}

	r.l.Debug("restoring job", "jobRow", jobRow)
	job, err := scripts.RestoreJob(
		jobRow.ID,
		jobRow.OwnerID,
		jobRow.ScriptID,
		jobRow.State,
		inputValues,
		outputs,
		"",
		result,
		jobRow.StartedAt,
		jobRow.ClosedAt,
		nil,
	)
	r.l.Debug("restored job", "job", *job, "err", err)
	if err != nil {
		r.l.Error("failed to restore job", "err", err.Error())
		return nil, err
	}

	r.l.Debug("returning job", "job", *job)
	return job, nil
}

const updateJobQuery = `
	UPDATE jobs
	SET 
		user_id = :user_id,
		script_id = :script_id,
		started_at = :started_at,
		closed_at = :closed_at,
		state = :state,
		status_code = :status_code,
		error_message = :error_message
	WHERE job_id = :job_id
`

const deleteFieldsQuery = `DELETE FROM job_params jp
USING parameters p
JOIN fields f ON p.field_id = f.field_id
WHERE jp.parameter_id = p.parameter_id
  AND jp.job_id = $1
  AND f.param = $2;
`

func (r *JobRepo) Update(ctx context.Context, job *scripts.Job) (err error) {
	r.l.Debug("updating job", "job", *job)

	r.l.Debug("beginning transaction")
	tx, err := r.db.BeginTxx(ctx, nil)
	r.l.Debug("transaction started", "err", err)
	if err != nil {
		r.l.Error("failed to start transaction", "err", err.Error())
		return err
	}
	defer func() {
		r.l.Debug("transaction finished", "err", err)
		if err != nil {
			r.l.Error("failed to commit transaction", "err", err.Error())
			_ = tx.Rollback()
		} else {
			r.l.Debug("transaction committed")
			err = tx.Commit()
		}
	}()

	r.l.Debug("converting job to db row", "job", *job)
	row := convertJobWithResToDB(job)
	r.l.Debug("converted job to db row", "row", row)

	r.l.Debug("updating job", "job", *job)
	res, err := tx.NamedExecContext(ctx, updateJobQuery, row)
	r.l.Debug("updated job", "res", res, "err", err)
	if err != nil {
		r.l.Error("failed to update job", "err", err.Error())
		return err
	}

	r.l.Debug("checking rows affected")
	rowsAffected, err := res.RowsAffected()
	r.l.Debug("rows affected", "rowsAffected", rowsAffected, "err", err)
	if err != nil {
		r.l.Error("failed to check rows affected", "err", err.Error())
		return err
	}
	r.l.Debug("сhecking if no rows affected", "is", rowsAffected == 0)
	if rowsAffected == 0 {
		r.l.Error("failed to update job", "jobID", job.ID(), "err", err.Error())
		return fmt.Errorf("%w Update: cannot update job with id: %d", scripts.ErrJobNotFound, job.ID())
	}

	r.l.Debug("deleting fields", "jobID", job.ID(), "field", "in")
	_, err = tx.ExecContext(ctx, deleteFieldsQuery, job.ID(), "in")
	r.l.Debug("deleted fields", "err", err)
	if err != nil {
		r.l.Error("failed to delete fields", "err", err.Error())
		return err
	}

	r.l.Debug("inserting fields", "jobID", job.ID(), "field", "in")
	if err := insertValuesTx(ctx, tx, int64(job.ID()), int64(job.ScriptID()), job.Input(), "in"); err != nil {
		r.l.Error("failed to insert fields", "err", err.Error())
		return err
	}
	r.l.Debug("inserted fields", "err", err)

	r.l.Debug("gettingjob result")
	jobRes, err := job.Result()
	r.l.Debug("job result", "jobRes", jobRes, "err", err)
	if err != nil {
		if errors.Is(err, scripts.ErrJobIsNotFinished) {
			r.l.Info("job is not finished")
			return nil
		}
		r.l.Error("failed to get job result", "err", err.Error())
		return err
	}

	r.l.Debug("job result", "jobRes", jobRes)
	if jobRes == nil {
		r.l.Info("job result is nil")
		return nil
	}
	r.l.Debug("deleting fields", "jobID", job.ID(), "field", "out")
	_, err = tx.ExecContext(ctx, deleteFieldsQuery, job.ID(), "out")
	r.l.Debug("deleted fields", "err", err)
	if err != nil {
		r.l.Error("failed to delete fields", "err", err.Error())
		return err
	}

	r.l.Debug("inserting fields", "jobID", job.ID(), "field", "out")
	if err := insertValuesTx(ctx, tx, int64(job.ID()), int64(job.ScriptID()), jobRes.Output(), "out"); err != nil {
		r.l.Error("failed to insert fields", "err", err.Error())
		return err
	}
	r.l.Debug("inserted fields", "err", err)

	return nil
}

const insertQuery = `
		INSERT INTO parameters (field_id, value)
		VALUES (:field_id, :value)
		RETURNING parameter_id
	`

const userJobsQuery = `
		SELECT *
		FROM jobs
		WHERE user_id = $1
		ORDER BY started_at DESC
	`

const userJobsWithStateQuery = `
		SELECT *
		FROM jobs
		WHERE user_id = $1 AND state = $2
		ORDER BY started_at DESC
	`

func (r *JobRepo) UserJobs(ctx context.Context, userID scripts.UserID) ([]scripts.Job, error) {
	r.l.Debug("getting user jobs", "userID", userID, "ctx", ctx)
	var jobRows []JobRow
	r.l.Debug("user jobs query", "userID", userID, "ctx", ctx)
	if err := r.db.SelectContext(ctx, &jobRows, userJobsQuery, userID); err != nil {
		r.l.Error("failed to get user jobs", "err", err.Error())
		return nil, fmt.Errorf("UserJobs: %w", err)
	}
	r.l.Debug("user jobs", "jobRows count", len(jobRows))

	return r.buildJobsFromRows(ctx, jobRows)
}

func (r *JobRepo) UserJobsWithState(ctx context.Context, userID scripts.UserID, jobState scripts.JobState) ([]scripts.Job, error) {
	r.l.Debug("getting user jobs with state", "userID", userID, "jobState", jobState, "ctx", ctx)
	var jobRows []JobRow
	r.l.Debug("user jobs with state query")
	if err := r.db.SelectContext(ctx, &jobRows, userJobsWithStateQuery, userID, jobState.String()); err != nil {
		r.l.Error("failed to get user jobs with state", "err", err.Error())
		return nil, err
	}

	r.l.Debug("user jobs with state", "jobRows count", len(jobRows))
	return r.buildJobsFromRows(ctx, jobRows)
}

func (r *JobRepo) buildJobsFromRows(ctx context.Context, rows []JobRow) ([]scripts.Job, error) {
	r.l.Debug("building jobs from rows", "rows count", len(rows), "ctx", ctx)
	var jobs []scripts.Job

	for _, jr := range rows {
		r.l.Debug("building job from row", "row", jr)
		inputValues, outputValues, err := getJobValues(ctx, r.db, jr.ID)
		r.l.Debug("got job values", "err", err)
		if err != nil {
			r.l.Error("failed to get job values", "err", err.Error())
			return nil, fmt.Errorf("buildJobsFromRows: getJobValues: %w", err)
		}

		var result *scripts.Result
		r.l.Debug("creating result")
		r.l.Debug("is needed to restore", "is", jr.StatusCode != nil)
		if jr.StatusCode != nil {
			result = scripts.RestoreResult(outputValues, scripts.StatusCode(*jr.StatusCode), jr.ErrorMessage)
			r.l.Debug("creating result", "result", *result)
		}

		var outFields []fieldRow
		r.l.Debug("getting fields", "ctx", ctx)
		err = r.db.SelectContext(ctx, &outFields, getFieldsQuery, jr.ScriptID, "out")
		r.l.Debug("got fields", "err", err)
		if err != nil {
			r.l.Error("failed to get fields", "err", err.Error())
			return nil, fmt.Errorf("buildJobsFromRows: SelectContext: %w", err)
		}

		r.l.Debug("converting fields to domain", "fields", outFields)
		outputs, err := convertFieldRowsToDomain(outFields)
		r.l.Debug("converted fields to domain", "err", err)
		if err != nil {
			r.l.Error("failed to convert fields to domain", "err", err.Error())
			return nil, fmt.Errorf("buildJobsFromRows: convertFieldRowsToDomain: %w", err)
		}
		r.l.Debug("restoring job", "jobID", jr.ID)
		job, err := scripts.RestoreJob(
			jr.ID,
			jr.OwnerID,
			jr.ScriptID,
			jr.State,
			inputValues,
			outputs,
			"",
			result,
			jr.StartedAt,
			jr.ClosedAt,
			nil,
		)
		r.l.Debug("restored job", "job", *job)
		if err != nil {
			r.l.Error("failed to restore job", "err", err.Error())
			return nil, fmt.Errorf("buildJobsFromRows: RestoreJob: %w", err)
		}
		jobs = append(jobs, *job)
	}

	r.l.Debug("returning jobs", "jobs count", len(jobs))
	return jobs, nil
}

func getJobValues(ctx context.Context, exec sqlx.ExtContext, jobID int64) (inputValues, outputValues []scripts.Value, err error) {
	var paramRows []ValueRow
	if err = sqlx.SelectContext(ctx, exec, &paramRows, paramsJobQuery, jobID); err != nil {
		return nil, nil, err
	}

	for _, pr := range paramRows {
		val, err := scripts.NewValue(pr.FieldType, pr.Value)
		if err != nil {
			return nil, nil, err
		}
		switch pr.Param {
		case "in":
			inputValues = append(inputValues, val)
		case "out":
			outputValues = append(outputValues, val)
		}
	}
	return inputValues, outputValues, nil
}

const fieldsQuery = `
		SELECT f.field_id, f.field_type, f.param
		FROM script_fields sf
		JOIN fields f ON sf.field_id = f.field_id
		WHERE sf.script_id = $1 AND f.param = $2
	`

const linkQuery = `
		INSERT INTO job_params (job_id, parameter_id)
		VALUES ($1, $2)
	`

func insertValuesTx(ctx context.Context, tx *sqlx.Tx, jobID, scriptID int64, values []scripts.Value, param string) error {
	var fields []ValueRow
	if err := tx.SelectContext(ctx, &fields, fieldsQuery, scriptID, param); err != nil {
		return err
	}
	for i, val := range values {
		row := convertValueToDB(val, fields[i].FieldID, fields[i].FieldType, fields[i].Param)

		var parameterID int64
		stmt, err := tx.PrepareNamedContext(ctx, insertQuery)
		if err != nil {
			return err
		}

		if err := stmt.GetContext(ctx, &parameterID, row); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, linkQuery, jobID, parameterID); err != nil {
			return err
		}

		if err := stmt.Close(); err != nil {
			return err
		}
	}

	return nil
}

type JobRow struct {
	ID           int64      `db:"job_id"`
	OwnerID      int64      `db:"user_id"`
	ScriptID     int64      `db:"script_id"`
	StartedAt    time.Time  `db:"started_at"`
	State        string     `db:"state"`
	ClosedAt     *time.Time `db:"closed_at"`
	StatusCode   *int64     `db:"status_code"`
	ErrorMessage *string    `db:"error_message"`
}

func convertJobPrototypeToDB(j *scripts.JobPrototype) JobRow {
	return JobRow{
		OwnerID:   int64(j.OwnerID()),
		ScriptID:  int64(j.ScriptID()),
		StartedAt: j.CreatedAt(),
	}
}

func convertJobToDB(j *scripts.Job) JobRow {
	return JobRow{
		ID:        int64(j.ID()),
		OwnerID:   int64(j.OwnerID()),
		ScriptID:  int64(j.ScriptID()),
		StartedAt: j.CreatedAt(),
		State:     j.State().String(),
	}
}

func convertJobWithResToDB(j *scripts.Job) JobRow {
	res, err := j.Result()
	if err != nil || res == nil {
		return convertJobToDB(j)
	}
	finishedAt, err := j.FinishedAt()
	if err != nil || finishedAt == nil {
		return convertJobToDB(j)
	}
	state := j.State().String()
	if state == "" {
		return convertJobToDB(j)
	}
	code := int64(res.Code())
	return JobRow{
		ID:           int64(j.ID()),
		OwnerID:      int64(j.OwnerID()),
		ScriptID:     int64(j.ScriptID()),
		StartedAt:    j.CreatedAt(),
		State:        state,
		StatusCode:   &code,
		ErrorMessage: res.ErrorMessage(),
		ClosedAt:     finishedAt,
	}
}

type ValueRow struct {
	FieldID   int64  `db:"field_id"`
	FieldType string `db:"field_type"`
	Value     string `db:"value"`
	Param     string `db:"param"`
}

func convertValueToDB(v scripts.Value, fieldID int64, fieldType string, param string) ValueRow {
	return ValueRow{
		FieldID:   fieldID,
		FieldType: fieldType,
		Value:     v.String(),
		Param:     param,
	}
}
