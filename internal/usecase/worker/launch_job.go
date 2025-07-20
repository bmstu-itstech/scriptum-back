package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type JobLaunchUC struct {
	jobR     scripts.JobRepository
	launcher scripts.Launcher
	notifier scripts.Notifier
}

func NewJobLaunchUC(
	jobR scripts.JobRepository,
	launcher scripts.Launcher,
	notifier scripts.Notifier,
) (*JobLaunchUC, error) {

	if jobR == nil {
		return nil, scripts.ErrInvalidJobRepository
	}
	if launcher == nil {
		return nil, scripts.ErrInvalidLauncherService
	}
	if notifier == nil {
		return nil, scripts.ErrInvalidNotifierService
	}

	return &JobLaunchUC{
		jobR:     jobR,
		launcher: launcher,
		notifier: notifier,
	}, nil
}

func (l *JobLaunchUC) ProcessLaunchRequest(ctx context.Context, input scripts.LaunchRequest) error {
	job := input.Job()

	result, err := l.launcher.Launch(ctx, job, input.ScriptFields())
	if err != nil {
		return err
	}

	err = l.jobR.CloseJob(ctx, job.JobID(), &result)
	if err != nil {
		return err
	}

	if input.NeedToNotify() {
		err = l.notifier.Notify(ctx, result, input.UserEmail())
		if err != nil {
			return err
		}
	}

	return nil
}
