package usecase

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
		panic(scripts.ErrInvalidJobRepository)
	}
	if launcher == nil {
		panic(scripts.ErrInvalidLauncherService)
	}
	if notifier == nil {
		panic(scripts.ErrInvalidNotifierService)
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

	err = l.jobR.Update(ctx, job.JobID(), &result)
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
