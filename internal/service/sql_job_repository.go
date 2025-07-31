package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/jmoiron/sqlx"
)

type JobRepo struct {
	db *sqlx.DB
}

func NewJobRepository(db *sqlx.DB) *JobRepo {
	return &JobRepo{
		db: db,
	}
}

const createJobQuery = `
	INSERT INTO jobs (user_id, script_id)
	VALUES (:user_id, :script_id)
	RETURNING job_id
`

func (r *JobRepo) Create(ctx context.Context, job *scripts.JobPrototype) (*scripts.Job, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	var jobID int64

	named, args, err := sqlx.Named(createJobQuery, convertJobPrototipeToDB(job))
	if err != nil {
		return nil, err
	}
	query := tx.Rebind(named)

	err = tx.QueryRowContext(ctx, query, args...).Scan(&jobID)
	if err != nil {
		return nil, err
	}

	if err := insertValuesTx(ctx, tx, jobID, int64(job.ScriptID()), job.Input(), "in"); err != nil {
		return nil, err
	}

	return job.Build(scripts.JobID(jobID))
}

const deleteJobQuery = `DELETE FROM jobs WHERE job_id = $1`

func (r *JobRepo) Delete(ctx context.Context, jobID scripts.JobID) error {
	result, err := r.db.ExecContext(ctx, deleteJobQuery, jobID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w Delete: cannot delete job with id: %d", scripts.ErrJobNotFound, jobID)
	}

	return nil
}

const getJobQuery = "SELECT * FROM jobs WHERE job_id=$1"

const getURLQuery = `
	SELECT path
	FROM scripts
	WHERE script_id = $1
`

const paramsJobQuery = `
		SELECT p.field_id, f.field_type, p.value, f.param
		FROM job_params jp
		JOIN parameters p ON jp.parameter_id = p.parameter_id
		JOIN fields f ON p.field_id = f.field_id
		WHERE jp.job_id = $1
		ORDER BY p.parameter_id
	`

func (r *JobRepo) Job(ctx context.Context, jobID scripts.JobID) (*scripts.Job, error) {
	var jobRow JobRow
	err := r.db.GetContext(ctx, &jobRow, getJobQuery, jobID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w Job: cannot extract job with id: %d", scripts.ErrJobNotFound, jobID)
		}
		return nil, err
	}
	var paramRows []ValueRow
	if err := r.db.SelectContext(ctx, &paramRows, paramsJobQuery, jobID); err != nil {
		return nil, err
	}

	inputValues, outputValues, err := getJobValues(ctx, r.db, jobRow.ID)
	if err != nil {
		return nil, err
	}

	var result *scripts.Result
	if jobRow.StatusCode != nil && jobRow.ClosedAt != nil {
		result = scripts.RestoreResult(outputValues, scripts.StatusCode(*jobRow.StatusCode), jobRow.ErrorMessage)
	}

	var path string
	err = r.db.GetContext(ctx, &path, getURLQuery, jobRow.ScriptID)
	if err != nil {
		return nil, err
	}

	var outFields []fieldRow
	err = r.db.SelectContext(ctx, &outFields, getFieldsQuery, jobRow.ScriptID, "out")
	if err != nil {
		return nil, err
	}
	outputs, err := convertFieldRowsToDomain(outFields)
	if err != nil {
		return nil, err
	}

	job, err := scripts.RestoreJob(
		jobRow.ID,
		jobRow.OwnerID,
		jobRow.ScriptID,
		jobRow.State,
		inputValues,
		outputs,
		path,
		result,
		jobRow.StartedAt,
		jobRow.ClosedAt,
	)
	if err != nil {
		return nil, err
	}

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
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	row := convertJobWithResToDB(job)

	res, err := tx.NamedExecContext(ctx, updateJobQuery, row)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("%w Update: cannot update job with id: %d", scripts.ErrJobNotFound, job.ID())
	}

	_, err = tx.ExecContext(ctx, deleteFieldsQuery, job.ID(), "in")
	if err != nil {
		return err
	}

	if err := insertValuesTx(ctx, tx, int64(job.ID()), int64(job.ScriptID()), job.Input(), "in"); err != nil {
		return err
	}

	jobRes, err := job.Result()
	if err != nil {
		if errors.Is(err, scripts.ErrJobIsNotFinished) {
			return nil
		}
		return err
	}

	if jobRes == nil {
		return nil
	}
	_, err = tx.ExecContext(ctx, deleteFieldsQuery, job.ID(), "out")
	if err != nil {
		return err
	}

	if err := insertValuesTx(ctx, tx, int64(job.ID()), int64(job.ScriptID()), jobRes.Output(), "out"); err != nil {
		return err
	}

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
	`

const userJobsWithStateQuery = `
		SELECT *
		FROM jobs
		WHERE user_id = $1 AND state = $2
	`

func (r *JobRepo) UserJobs(ctx context.Context, userID scripts.UserID) ([]scripts.Job, error) {

	var jobRows []JobRow
	if err := r.db.SelectContext(ctx, &jobRows, userJobsQuery, userID); err != nil {
		return nil, err
	}

	return r.buildJobsFromRows(ctx, jobRows)
}

func (r *JobRepo) UserJobsWithState(ctx context.Context, userID scripts.UserID, jobState scripts.JobState) ([]scripts.Job, error) {

	var jobRows []JobRow
	if err := r.db.SelectContext(ctx, &jobRows, userJobsWithStateQuery, userID, jobState.String()); err != nil {
		return nil, err
	}

	return r.buildJobsFromRows(ctx, jobRows)
}

func (r *JobRepo) buildJobsFromRows(ctx context.Context, rows []JobRow) ([]scripts.Job, error) {
	var jobs []scripts.Job

	for _, jr := range rows {
		inputValues, outputValues, err := getJobValues(ctx, r.db, jr.ID)
		if err != nil {
			return nil, err
		}

		var result *scripts.Result
		if jr.StatusCode != nil {
			result = scripts.RestoreResult(outputValues, scripts.StatusCode(*jr.StatusCode), jr.ErrorMessage)
		}

		var path string
		err = r.db.GetContext(ctx, &path, getURLQuery, jr.ScriptID)
		if err != nil {
			return nil, err
		}

		var outFields []fieldRow
		err = r.db.SelectContext(ctx, &outFields, getFieldsQuery, jr.ScriptID, "out")
		if err != nil {
			return nil, err
		}
		outputs, err := convertFieldRowsToDomain(outFields)
		if err != nil {
			return nil, err
		}
		job, err := scripts.RestoreJob(
			jr.ID,
			jr.OwnerID,
			jr.ScriptID,
			jr.State,
			inputValues,
			outputs,
			path,
			result,
			jr.StartedAt,
			jr.ClosedAt,
		)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, *job)
	}

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
		defer stmt.Close()

		if err := stmt.GetContext(ctx, &parameterID, row); err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, linkQuery, jobID, parameterID); err != nil {
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

func convertJobPrototipeToDB(j *scripts.JobPrototype) JobRow {
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
