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
	u.logger.Debug("data", "userID", userID, "jobID", jobID)

	u.logger.Debug("get job from repository", "jobID", jobID)
	job, err := u.jobR.Job(ctx, scripts.JobID(jobID))
	u.logger.Debug("got job", "job", *job, "err", err.Error())
	if err != nil {
		u.logger.Error("failed to get job", "err", err.Error())
		return JobDTO{}, err
	}

	u.logger.Debug("check job owner", "job", *job, "userID", userID)
	u.logger.Debug("is owner", "is", job.OwnerID() == scripts.UserID(userID), "userID", userID, "jobID", jobID)

	if job.OwnerID() != scripts.UserID(userID) {
		u.logger.Error("failed to get job", "err", scripts.ErrPermissionDenied.Error())
		return JobDTO{}, scripts.ErrPermissionDenied
	}

	u.logger.Debug("get script from repository", "job", *job)
	script, err := u.scriptR.Script(ctx, job.ScriptID())
	u.logger.Debug("got script", "script", script, "err", err.Error())
	if err != nil {
		u.logger.Error("failed to get job", "err", err.Error())
		return JobDTO{}, err
	}

	u.logger.Info("convert job to dto")
	u.logger.Debug("convert job to dto", "job", *job, "script name", script.Name())
	return JobToDTO(*job, script.Name())
}
