package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/zhikh23/pgutils"

	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

func (r *Repository) Job(ctx context.Context, id value.JobID) (*entity.Job, error) {
	l := r.l.With(
		slog.String("op", "postgres.Repository.Job"),
		slog.String("job_id", string(id)),
	)

	var job *entity.Job
	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		var err error
		job, err = r.job(ctx, tx, id)
		return err
	})
	if errors.Is(err, sql.ErrNoRows) {
		l.WarnContext(ctx, "job not found", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %s", ports.ErrJobNotFound, err.Error())
	}
	if err != nil {
		l.ErrorContext(ctx, "failed to execute transaction", slog.String("error", err.Error()))
		return nil, err
	}
	return job, nil
}

func (r *Repository) job(ctx context.Context, qc sqlx.QueryerContext, id value.JobID) (*entity.Job, error) {
	rJob, err := r.selectJobRow(ctx, qc, string(id))
	if err != nil {
		return nil, err
	}
	rInput, err := r.selectJobInputValueRows(ctx, qc, string(id))
	if err != nil {
		return nil, err
	}
	rOutput, err := r.selectJobOutputValueRows(ctx, qc, string(id))
	if err != nil {
		return nil, err
	}
	rOut, err := r.selectJobOutputFieldRows(ctx, qc, string(id))
	if err != nil {
		return nil, err
	}
	job, err := jobRowToDomain(rJob, rInput, rOutput, rOut)
	if err != nil {
		return nil, err
	}
	return job, nil
}

func (r *Repository) UserJobs(ctx context.Context, uid value.UserID) ([]*entity.Job, error) {
	l := r.l.With(
		slog.String("op", "postgres.Repository.UserJobs"),
		slog.Int64("user_id", int64(uid)),
	)

	var jobs []*entity.Job
	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		var err error
		jobs, err = r.userJobs(ctx, tx, uid)
		return err
	})
	if err != nil {
		l.ErrorContext(ctx, "failed to execute transaction", slog.String("error", err.Error()))
		return nil, err
	}
	return jobs, nil
}

func (r *Repository) userJobs(ctx context.Context, qc sqlx.QueryerContext, uid value.UserID) ([]*entity.Job, error) {
	jobs := make([]*entity.Job, 0)
	rJobs, err := r.selectUserJobRows(ctx, qc, int64(uid))
	if err != nil {
		return nil, err
	}
	for _, rJob := range rJobs {
		rInput, err2 := r.selectJobInputValueRows(ctx, qc, rJob.ID)
		if err2 != nil {
			return nil, err2
		}
		rOutput, err2 := r.selectJobOutputValueRows(ctx, qc, rJob.ID)
		if err2 != nil {
			return nil, err2
		}
		rOut, err2 := r.selectJobOutputFieldRows(ctx, qc, rJob.ID)
		if err2 != nil {
			return nil, err2
		}
		job, err2 := jobRowToDomain(rJob, rInput, rOutput, rOut)
		if err2 != nil {
			return nil, err2
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}

func (r *Repository) UserJobsWithState(
	ctx context.Context,
	uid value.UserID,
	state value.JobState,
) ([]*entity.Job, error) {
	l := r.l.With(
		slog.String("op", "postgres.Repository.UserJobsWithState"),
		slog.Int64("user_id", int64(uid)),
		slog.String("state", state.String()),
	)

	var jobs []*entity.Job
	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		var err error
		jobs, err = r.userJobsWithState(ctx, tx, uid, state)
		return err
	})
	if err != nil {
		l.ErrorContext(ctx, "failed to execute transaction", slog.String("error", err.Error()))
		return nil, err
	}
	return jobs, nil
}

func (r *Repository) userJobsWithState(
	ctx context.Context,
	qc sqlx.QueryerContext,
	uid value.UserID,
	state value.JobState,
) ([]*entity.Job, error) {
	jobs := make([]*entity.Job, 0)
	rJobs, err := r.selectUserJobRowsWithState(ctx, qc, int64(uid), state.String())
	if err != nil {
		return nil, err
	}
	for _, rJob := range rJobs {
		rInput, err2 := r.selectJobInputValueRows(ctx, qc, rJob.ID)
		if err2 != nil {
			return nil, err2
		}
		rOutput, err2 := r.selectJobOutputValueRows(ctx, qc, rJob.ID)
		if err2 != nil {
			return nil, err2
		}
		rOut, err2 := r.selectJobOutputFieldRows(ctx, qc, rJob.ID)
		if err2 != nil {
			return nil, err2
		}
		job, err2 := jobRowToDomain(rJob, rInput, rOutput, rOut)
		if err2 != nil {
			return nil, err2
		}
		jobs = append(jobs, job)
	}
	return jobs, nil
}
