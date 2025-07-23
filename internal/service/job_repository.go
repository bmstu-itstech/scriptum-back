package service

import (
	"context"
	"fmt"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/jackc/pgx/v4"
)

type JobRepo struct {
	DB SQLDBConn
}

func NewJobRepo(db SQLDBConn) *JobRepo {
	return &JobRepo{
		DB: db,
	}
}

const PostJobQuery = `
		INSERT INTO jobs (user_id, script_id, started_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		RETURNING job_id;
	`

func (r *JobRepo) Post(ctx context.Context, job scripts.Job, scriptID scripts.ScriptID) (scripts.JobID, error) {
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

func (r *JobRepo) Update(ctx context.Context, jobID scripts.JobID, res *scripts.Result) error {
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
