package worker

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type LaunchUC struct {
	jobR     scripts.JobRepository
	launcher PythonLauncher
	notifier scripts.Notifier
}

func NewLaunchUC(
	jobR scripts.JobRepository,
	launcher PythonLauncher,
	notifier scripts.Notifier,
) (*LaunchUC, error) {

	if jobR == nil {
		return nil, scripts.ErrInvalidJobRepository
	}
	// if launcher == nil {
	// 	return nil, scripts.ErrInvalidLauncherService
	// }
	if notifier == nil {
		return nil, scripts.ErrInvalidNotifierService
	}

	return &LaunchUC{
		jobR:     jobR,
		launcher: launcher,
		notifier: notifier,
	}, nil
}

func (l *LaunchUC) ProcessLaunchRequest(ctx context.Context, input scripts.LaunchRequest) error {
	// надо запустить скрипт, опубликовать сообщение

	result, err := l.launcher.Launch(ctx, input.Job, input.ScriptFields)
	if err != nil {
		return err
	}

	err = l.jobR.CloseJob(ctx, input.Job.JobID(), &result)
	if err != nil {
		return err
	}

	if input.NeedToNotify {
		err = l.notifier.Notify(ctx, result, input.GetUserEmail())
		if err != nil {
			return err
		}
	}
	return nil
}
