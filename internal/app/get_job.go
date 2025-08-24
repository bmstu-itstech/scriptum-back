package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetJobUC struct {
	jobR    scripts.JobRepository
	userP   scripts.UserProvider
	scriptR scripts.ScriptRepository
	logger  *slog.Logger
}

func NewGetJobUC(
	jobR scripts.JobRepository,
	userP scripts.UserProvider,
	scriptR scripts.ScriptRepository,
	logger *slog.Logger,
) GetJobUC {
	return GetJobUC{jobR: jobR, userP: userP, scriptR: scriptR, logger: logger}
}

func (u *GetJobUC) Job(ctx context.Context, userID uint32, jobID int64) (JobDTO, error) {
	u.logger.Info("get job", "jobID", jobID)

	job, err := u.jobR.Job(ctx, scripts.JobID(jobID))
	if err != nil {
		u.logger.Error("failed to get job", "err", err)
		return JobDTO{}, err
	}

	if job.OwnerID() != scripts.UserID(userID) {
		u.logger.Error("failed to get job", "err", scripts.ErrPermissionDenied)
		return JobDTO{}, scripts.ErrPermissionDenied
	}

	script, err := u.scriptR.Script(ctx, job.ScriptID())
	if err != nil {
		u.logger.Error("failed to get job", "err", err)
		return JobDTO{}, err
	}

	return JobToDTO(*job, script.Name())
}
