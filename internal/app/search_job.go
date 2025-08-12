package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type SearchJobsUC struct {
	jobR   scripts.JobRepository
	userP  scripts.UserProvider
	logger *slog.Logger
}

func NewSearchJobsUC(
	jobR scripts.JobRepository,
	userP scripts.UserProvider,
	logger *slog.Logger,
) SearchJobsUC {
	return SearchJobsUC{jobR: jobR, userP: userP, logger: logger}
}

func (u *SearchJobsUC) Search(ctx context.Context, userID uint32, state string) ([]JobDTO, error) {
	u.logger.Info("searching jobs ", "userID", userID)
	_, err := u.userP.User(ctx, scripts.UserID(userID))
	if err != nil {
		u.logger.Error("failed to search job", "err", err)
		return nil, err
	}

	jobState, err := scripts.NewJobStateFromString(state)
	if err != nil {
		u.logger.Error("failed to search job", "err", err)
		return nil, err
	}

	jobs, err := u.jobR.UserJobsWithState(ctx, scripts.UserID(userID), jobState)
	if err != nil {
		u.logger.Error("failed to search job", "err", err)
		return nil, err
	}

	dto := make([]JobDTO, 0, len(jobs))
	for _, j := range jobs {
		job, err := JobToDTO(j)
		if err != nil {
			u.logger.Error("failed to search job", "err", err)
			return nil, err
		}
		dto = append(dto, job)
	}

	return dto, nil
}
