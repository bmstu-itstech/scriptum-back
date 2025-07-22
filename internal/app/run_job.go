package app

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type JobRunUC struct {
	jobR     scripts.JobRepository
	launcher scripts.Launcher
	notifier scripts.Notifier
}

func NewJobRunUC(
	jobR scripts.JobRepository,
	launcher scripts.Launcher,
	notifier scripts.Notifier,
) (*JobRunUC, error) {

	if jobR == nil {
		return nil, scripts.ErrInvalidJobRepository
	}
	if launcher == nil {
		return nil, scripts.ErrInvalidLauncherService
	}
	if notifier == nil {
		return nil, scripts.ErrInvalidNotifierService
	}

	return &JobRunUC{
		jobR:     jobR,
		launcher: launcher,
		notifier: notifier,
	}, nil
}

func (l *JobRunUC) ProcessLaunchRequest(ctx context.Context, jobDTO JobDTO) error {
	job, err := DTOToJob(jobDTO)
	if err != nil {
		return err
	}

	result, err := l.launcher.Launch(ctx, job)
	if err != nil {
		return err
	}

	err = l.jobR.CloseJob(ctx, job.JobID(), &result)
	if err != nil {
		return err
	}

	if job.NeedToNotify() {
		err = l.notifier.Notify(ctx, result, job.UserEmail())
		if err != nil {
			return err
		}
	}

	return nil
}
