package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type JobRunUC struct {
	scriptR  scripts.ScriptRepository
	jobR     scripts.JobRepository
	runner   scripts.Runner
	notifier scripts.Notifier
	userP    scripts.UserProvider
	logger   *slog.Logger
}

func NewJobRunUC(
	scriptR scripts.ScriptRepository,
	jobR scripts.JobRepository,
	launcher scripts.Runner,
	notifier scripts.Notifier,
	userP scripts.UserProvider,
	logger *slog.Logger,
) JobRunUC {
	return JobRunUC{
		scriptR:  scriptR,
		jobR:     jobR,
		runner:   launcher,
		notifier: notifier,
		userP:    userP,
		logger:   logger,
	}
}

func (l *JobRunUC) Run(ctx context.Context, req JobDTO) error {
	l.logger.Info("running job", "req", req)
	job, err := DTOToJob(req)
	if err != nil {
		l.logger.Error("failed to run job", "err", err)
		return err
	}

	err = job.Run()
	if err != nil {
		l.logger.Error("failed to run job", "err", err)
		return err
	}

	err = l.jobR.Update(ctx, job)
	if err != nil {
		l.logger.Error("failed to run job", "err", err)
		return err
	}

	res, err := l.runner.Run(ctx, job)
	if err != nil {
		l.logger.Error("failed to run job", "err", err)
		return err
	}

	err = job.Finish(res)
	if err != nil {
		l.logger.Error("failed to run job", "err", err)
		return err
	}

	err = l.jobR.Update(ctx, job)
	if err != nil {
		l.logger.Error("failed to run job", "err", err)
		return err
	}

	user, err := l.userP.User(ctx, job.OwnerID())
	if err != nil {
		l.logger.Error("failed to run job", "err", err)
		return err
	}

	if req.NeedToNotify {
		err = l.notifier.Notify(ctx, job, user.Email())
		if err != nil {
			l.logger.Error("failed to run job", "err", err)
			return err
		}
	}

	return nil
}
