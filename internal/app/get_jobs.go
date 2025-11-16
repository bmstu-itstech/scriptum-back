package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetJobsUC struct {
	jobR    scripts.JobRepository
	userP   scripts.UserProvider
	scriptR scripts.ScriptRepository
	logger  *slog.Logger
}

func NewGetJobsUC(
	jobR scripts.JobRepository,
	userP scripts.UserProvider,
	scriptR scripts.ScriptRepository,
	logger *slog.Logger,
) GetJobsUC {
	return GetJobsUC{jobR: jobR, userP: userP, scriptR: scriptR, logger: logger}
}

func (u *GetJobsUC) Jobs(ctx context.Context, userID uint32) ([]JobDTO, error) {
	u.logger.Info("get jobs for user", "userID", userID)

	u.logger.Debug("get jobs for user", "userID", userID)
	jobs, err := u.jobR.UserJobs(ctx, scripts.UserID(userID))
	u.logger.Debug("got jobs", "jobs count", len(jobs), "err", err)
	if err != nil {
		u.logger.Error("failed to get jobs for user", "err", err.Error())
		return nil, err
	}

	dto := make([]JobDTO, 0, len(jobs))
	for _, j := range jobs {
		u.logger.Debug("get script for job", "jobID", j.ID())
		script, err := u.scriptR.Script(ctx, j.ScriptID())
		u.logger.Debug("got script", "script", script, "err", err)
		if err != nil {
			u.logger.Error("failed to get jobs for user", "err", err.Error())
			return nil, err
		}

		u.logger.Debug("convert job to dto", "job", j, "script name", script.Name())
		job, err := JobToDTO(j, script.Name())
		u.logger.Debug("converted job to dto", "job", job, "err", err)
		if err != nil {
			u.logger.Error("failed to get jobs for user", "err", err.Error())
			return nil, err
		}
		dto = append(dto, job)
	}

	u.logger.Info("got jobs for user", "userID", userID, "jobs count", len(dto))
	return dto, nil
}
