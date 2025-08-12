package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetJobsUC struct {
	jobR   scripts.JobRepository
	userP  scripts.UserProvider
	logger *slog.Logger
}

func NewGetJobsUC(
	jobR scripts.JobRepository,
	userP scripts.UserProvider,
	logger *slog.Logger,
) GetJobsUC {
	return GetJobsUC{jobR: jobR, userP: userP, logger: logger}
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
		job, err := JobToDTO(j)
		if err != nil {
			u.logger.Error("failed to get jobs for user", "err", err)
			return nil, err
		}
		dto = append(dto, job)
	}

	return dto, nil
}
