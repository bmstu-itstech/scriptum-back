package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type SearchJobsUC struct {
	jobR    scripts.JobRepository
	userP   scripts.UserProvider
	scriptR scripts.ScriptRepository
	logger  *slog.Logger
}

func NewSearchJobsUC(
	jobR scripts.JobRepository,
	userP scripts.UserProvider,
	scriptR scripts.ScriptRepository,
	logger *slog.Logger,
) SearchJobsUC {
	return SearchJobsUC{jobR: jobR, userP: userP, scriptR: scriptR, logger: logger}
}

func (u *SearchJobsUC) Search(ctx context.Context, userID uint32, state string) ([]JobDTO, error) {
	u.logger.Info("searching jobs", "userID", userID)
	u.logger.Debug("get user", "userID", userID, "state", state)
	_, err := u.userP.User(ctx, scripts.UserID(userID))
	u.logger.Debug("got user", "err", err.Error())
	if err != nil {
		u.logger.Error("failed to search job", "err", err.Error())
		return nil, err
	}

	u.logger.Debug("forming job state during to domain", "state", state)
	jobState, err := scripts.NewJobStateFromString(state)
	if err != nil {
		u.logger.Error("failed to search job", "err", err.Error())
		return nil, err
	}

	u.logger.Debug("searching jobs with state", "userID", userID, "state", jobState)
	jobs, err := u.jobR.UserJobsWithState(ctx, scripts.UserID(userID), jobState)
	u.logger.Debug("got jobs", "jobs count", len(jobs), "err", err.Error())
	if err != nil {
		u.logger.Error("failed to search job", "err", err.Error())
		return nil, err
	}

	u.logger.Debug("converting jobs to DTOs", "jobs count", len(jobs))
	dto := make([]JobDTO, 0, len(jobs))
	for _, j := range jobs {
		u.logger.Debug("getting script", "scriptID", j.ScriptID())
		script, err := u.scriptR.Script(ctx, j.ScriptID())
		u.logger.Debug("got script", "err", err.Error())
		if err != nil {
			u.logger.Error("failed to search job", "err", err.Error())
			return nil, err
		}

		u.logger.Debug("converting job to DTO", "jobID", j.ID())
		job, err := JobToDTO(j, script.Name())
		u.logger.Debug("got job DTO", "err", err.Error())
		if err != nil {
			u.logger.Error("failed to search job", "err", err.Error())
			return nil, err
		}
		dto = append(dto, job)
	}

	u.logger.Info("returning jobs", "jobs count", len(dto))
	return dto, nil
}
