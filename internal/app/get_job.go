package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetJobUC struct {
	jobR   scripts.JobRepository
	userP  scripts.UserProvider
	logger *slog.Logger
}

func NewGetJobUC(
	jobR scripts.JobRepository,
	userP scripts.UserProvider,
	logger *slog.Logger,
) GetJobUC {
	return GetJobUC{jobR: jobR, userP: userP, logger: logger}
}

func (u *GetJobUC) Job(ctx context.Context, userID uint32, jobID int64) (JobDTO, error) {
	u.logger.Info("get job", "jobID", jobID)
	_, err := u.userP.User(ctx, scripts.UserID(userID))
	if err != nil {
		u.logger.Error("failed to get job", "err", err)
		return JobDTO{}, err
	}

	job, err := u.jobR.Job(ctx, scripts.JobID(jobID))
	if err != nil {
		u.logger.Error("failed to get job", "err", err)
		return JobDTO{}, err
	}

	if job.OwnerID() != scripts.UserID(userID) {
		u.logger.Error("failed to get job", "err", scripts.ErrPermissionDenied)
		return JobDTO{}, scripts.ErrPermissionDenied
	}
	return JobToDTO(*job)
}
