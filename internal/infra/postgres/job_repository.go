package postgres

import (
	"context"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/zhikh23/pgutils"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

func (r *Repository) SaveJob(ctx context.Context, job *entity.Job) error {
	l := r.l.With(
		slog.String("op", "postgres.Repository.SaveJob"),
		slog.String("job_id", string(job.ID())),
	)
	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		return r.saveJob(ctx, tx, job)
	})
	if err != nil {
		l.ErrorContext(ctx, "failed to execute transaction", slog.String("error", err.Error()))
		return err
	}
	return nil
}

func (r *Repository) saveJob(ctx context.Context, ec sqlx.ExecerContext, job *entity.Job) error {
	rJob := jobRowFromDomain(job)
	if err := r.insertJobRow(ctx, ec, rJob); err != nil {
		return err
	}
	if len(job.Input()) > 0 {
		rInput := jobValueRowsFromDomain(job.Input(), job.ID())
		if err := r.insertJobInputValueRows(ctx, ec, rInput); err != nil {
			return err
		}
	}
	if res := job.Result(); res != nil {
		if len(res.Output()) > 0 {
			rOutput := jobValueRowsFromDomain(res.Output(), job.ID())
			if err := r.insertJobOutputValueRows(ctx, ec, rOutput); err != nil {
				return err
			}
		}
	}
	if len(job.Out()) > 0 {
		rOut := jobFieldRowsFromDomain(job.Out(), job.ID())
		if err := r.insertJobOutputFieldRows(ctx, ec, rOut); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repository) UpdateJob(
	ctx context.Context,
	id value.JobID,
	updateFn func(ctx2 context.Context, job *entity.Job) error,
) error {
	l := r.l.With(
		slog.String("op", "postgres.Repository.UpdateJob"),
		slog.String("job_id", string(id)),
	)
	err := pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		job, err := r.job(ctx, tx, id)
		if err != nil {
			return err
		}
		err = updateFn(ctx, job)
		if err != nil {
			return err
		}
		return r.saveJob(ctx, tx, job)
	})
	if err != nil {
		l.ErrorContext(ctx, "failed to execute transaction", slog.String("error", err.Error()))
		return err
	}
	return nil
}
