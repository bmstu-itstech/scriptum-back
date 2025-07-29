package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type JobRunUC struct {
	jobR     scripts.JobRepository
	runner   scripts.Runner
	notifier scripts.Notifier
	userP    scripts.UserProvider
	logger   *slog.Logger
}

func NewJobRunUC(
	jobR scripts.JobRepository,
	launcher scripts.Runner,
	notifier scripts.Notifier,
	userP scripts.UserProvider,
	logger *slog.Logger,
) JobRunUC {
	return JobRunUC{
		jobR:     jobR,
		runner:   launcher,
		notifier: notifier,
		userP:    userP,
		logger:   logger,
	}
}

func (l *JobRunUC) Run(ctx context.Context, jobDTO JobDTO, needToNotify bool) error {
	job, err := DTOToJob(jobDTO)
	if err != nil {
		return err
	}

	err = job.Run()
	if err != nil {
		return err
	}

	err = l.jobR.Update(ctx, job)
	if err != nil {
		return err
	}

	res, err := l.runner.Run(ctx, job)
	if err != nil {
		return err
	}

	err = job.Finish(res)
	if err != nil {
		return err
	}

	err = l.jobR.Update(ctx, job)
	if err != nil {
		return err
	}

	user, err := l.userP.User(ctx, job.OwnerID())
	if err != nil {
		return err
	}

	if needToNotify {
		err = l.notifier.Notify(ctx, job, user.Email())
		if err != nil {
			return err
		}
	}

	return nil
}
