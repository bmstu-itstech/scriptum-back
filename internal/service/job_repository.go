package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/jackc/pgx/v4"
)

type JobRepo struct {
	DB SQLDBConn
}

func NewJobRepo(ctx context.Context) (*JobRepo, error) {
	host := "localhost"
	port := 5432
	user := "app_user"
	password := "your_secure_password"
	dbname := "dev"

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, password, host, port, dbname)

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}

	return &JobRepo{
		DB: conn,
	}, nil
}

const PostJobQuery = `
		INSERT INTO jobs (user_id, script_id, started_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		RETURNING job_id;
	`

func (r *JobRepo) PostJob(ctx context.Context, job scripts.Job, scriptID scripts.ScriptID) (scripts.JobID, error) {
	var rawID int
	err := r.DB.QueryRow(ctx, PostJobQuery,
		job.UserID(),
		scriptID,
	).Scan(&rawID)
	if err != nil {
		return 0, err
	}
	return scripts.JobID(rawID), nil
}

const CloseJobQuery = `
	UPDATE jobs SET
		status_code = $1,
		error_message = $2,
		closed_at = CURRENT_TIMESTAMP
	WHERE job_id = $3;
`

const insertOutParamQuery = `
		INSERT INTO parameters (value) VALUES ($1) RETURNING parameter_id;
	`

const insertJobParamQuery = `
		INSERT INTO job_params (job_id, parameter_id, param)
		VALUES ($1, $2, 'out');
	`

func (r *JobRepo) CloseJob(ctx context.Context, jobID scripts.JobID, res *scripts.Result) error {
	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			fmt.Printf("tx rollback error: %v", err)
		}
	}()

	_, err = tx.Exec(ctx, CloseJobQuery,
		res.Code(),
		*res.ErrorMessage(),
		jobID,
	)

	if err != nil {
		return err
	}

	for _, val := range res.Out().Get() {
		var paramID int64
		err := tx.QueryRow(ctx, insertOutParamQuery, val).Scan(&paramID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, insertJobParamQuery, jobID, paramID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

const JobsByScriptIDQuery = `
	SELECT
		j.job_id,
		j.user_id,
		j.started_at,
		j.script_id,
		p.value,
		jp.param,
		f.field_type
	FROM jobs j
	LEFT JOIN job_params jp ON j.job_id = jp.job_id
	LEFT JOIN parameters p ON p.parameter_id = jp.parameter_id
	LEFT JOIN fields f ON f.field_id = p.field_id
	WHERE j.script_id = $1;
`

type JobAccumulator struct {
	userID    scripts.UserID
	startedAt time.Time
	inputVals []scripts.Value
}

func (r *JobRepo) JobsByScriptID(ctx context.Context, scriptID scripts.ScriptID) ([]scripts.Job, error) {
	rows, err := r.DB.Query(ctx, JobsByScriptIDQuery, scriptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobMap := make(map[scripts.JobID]*JobAccumulator)

	for rows.Next() {
		var (
			jobID     int
			userID    scripts.UserID
			startedAt time.Time

			valStr    *string
			paramType *string
			fieldType string
		)

		if err := rows.Scan(&jobID, &userID, &startedAt, new(int64), &valStr, &paramType, &fieldType); err != nil {
			return nil, err
		}

		jid := scripts.JobID(jobID)
		acc, ok := jobMap[jid]
		if !ok {
			acc = &JobAccumulator{
				userID:    userID,
				startedAt: startedAt,
			}
			jobMap[jid] = acc
		}

		if valStr != nil && paramType != nil && *paramType == "in" {
			val, err := scripts.NewValue(fieldType, *valStr)
			if err != nil {
				return nil, err
			}
			acc.inputVals = append(acc.inputVals, val)
		}
	}

	var jobs []scripts.Job
	for jid, acc := range jobMap {
		inVec, err := scripts.NewVector(acc.inputVals)
		if err != nil {
			return nil, err
		}
		job, err := scripts.NewJob(jid, acc.userID, *inVec, "", acc.startedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, *job)
	}

	return jobs, nil
}
