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
	u.logger.Info("searching jobs ", "userID", userID)
	_, err := u.userP.User(ctx, scripts.UserID(userID))
	if err != nil {
		u.logger.Error("failed to search job", "err", err.Error())
		return nil, err
	}

	jobState, err := scripts.NewJobStateFromString(state)
	if err != nil {
		u.logger.Error("failed to search job", "err", err.Error())
		return nil, err
	}

	jobs, err := u.jobR.UserJobsWithState(ctx, scripts.UserID(userID), jobState)
	if err != nil {
		u.logger.Error("failed to search job", "err", err.Error())
		return nil, err
	}

	dto := make([]JobDTO, 0, len(jobs))
	for _, j := range jobs {
		script, err := u.scriptR.Script(ctx, j.ScriptID())
		if err != nil {
			u.logger.Error("failed to search job", "err", err.Error())
			return nil, err
		}

		job, err := JobToDTO(j, script.Name())
		if err != nil {
			u.logger.Error("failed to search job", "err", err.Error())
			return nil, err
		}
		dto = append(dto, job)
	}

	return dto, nil
}
