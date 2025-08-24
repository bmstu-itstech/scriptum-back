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

	jobs, err := u.jobR.UserJobs(ctx, scripts.UserID(userID))
	if err != nil {
		u.logger.Error("failed to get jobs for user", "err", err)
		return nil, err
	}

	dto := make([]JobDTO, 0, len(jobs))
	for _, j := range jobs {
		script, err := u.scriptR.Script(ctx, j.ScriptID())
		if err != nil {
			u.logger.Error("failed to get jobs for user", "err", err)
			return nil, err
		}

		job, err := JobToDTO(j, script.Name())
		if err != nil {
			u.logger.Error("failed to get jobs for user", "err", err)
			return nil, err
		}
		dto = append(dto, job)
	}

	return dto, nil
}
