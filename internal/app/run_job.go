package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type JobRunUC struct {
	scriptR  scripts.ScriptRepository
	jobR     scripts.JobRepository
	runner   scripts.Launcher
	notifier scripts.Notifier
	userP    scripts.UserProvider
	manager  scripts.FileManager
	logger   *slog.Logger
}

func NewJobRunUC(
	scriptR scripts.ScriptRepository,
	jobR scripts.JobRepository,
	launcher scripts.Launcher,
	notifier scripts.Notifier,
	userP scripts.UserProvider,
	manager scripts.FileManager,
	logger *slog.Logger,
) JobRunUC {
	return JobRunUC{
		scriptR:  scriptR,
		jobR:     jobR,
		runner:   launcher,
		notifier: notifier,
		userP:    userP,
		manager:  manager,
		logger:   logger,
	}
}

func (l *JobRunUC) Run(ctx context.Context, req JobDTO) error {
	l.logger.Info("running job", "req", req)
	l.logger.Debug("converting dto to job", "req", req, "ctx", ctx)
	job, err := DTOToJob(req)
	l.logger.Debug("converted job", "job", job, "err", err.Error())
	if err != nil {
		l.logger.Error("failed to run job", "err", err.Error())
		return err
	}

	l.logger.Debug("running job", "job", job)
	err = job.Run()
	l.logger.Debug("job finished", "err", err.Error())
	if err != nil {
		l.logger.Error("failed to run job", "err", err.Error())
		return err
	}

	l.logger.Debug("updating job to \"running\"", "job", job)
	err = l.jobR.Update(ctx, job)
	l.logger.Debug("updated job", "err", err.Error())
	if err != nil {
		l.logger.Error("failed to run job", "err", err.Error())
		return err
	}

	l.logger.Debug("running job with Runner", "job", job)
	res, err := l.runner.Run(ctx, job)
	l.logger.Debug("runner finished", "res", res, "err", err.Error())
	if err != nil {
		l.logger.Error("failed to run job", "err", err.Error())
		return err
	}

	l.logger.Info("job finished", "res", res)

	l.logger.Debug("finishing job", "job", job)
	err = job.Finish(res)
	l.logger.Debug("job finished", "err", err.Error())
	if err != nil {
		l.logger.Error("failed to run job", "err", err.Error())
		return err
	}

	l.logger.Debug("updating job to \"finished\"", "job", job)
	err = l.jobR.Update(ctx, job)
	l.logger.Debug("updated job", "err", err.Error())
	if err != nil {
		l.logger.Error("failed to run job", "err", err.Error())
		return err
	}

	l.logger.Debug("getting user", "job", job)
	user, err := l.userP.User(ctx, job.OwnerID())
	l.logger.Debug("got user", "user", user, "err", err.Error())
	if err != nil {
		l.logger.Error("failed to run job", "err", err.Error())
		return err
	}

	l.logger.Debug("notifying user", "needed", req.NeedToNotify)
	if req.NeedToNotify {
		l.logger.Debug("notifying user", "user mail", user.Email(), "job", job)
		err = l.notifier.Notify(ctx, job, user.Email())
		l.logger.Debug("notified user", "err", err.Error())
		if err != nil {
			l.logger.Error("failed to run job", "err", err.Error())
			return err
		}
	}

	l.logger.Debug("deleting sandbox", "url", req.Url)
	err = l.runner.DeleteSandbox(ctx, req.Url)
	l.logger.Debug("deleted sandbox", "err", err.Error())
	if err != nil {
		l.logger.Error("failed to delete sandbox", "err", err.Error())
		return err
	}

	return nil
}
