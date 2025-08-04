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
	_, err := u.userP.User(ctx, scripts.UserID(userID))
	if err != nil {
		return JobDTO{}, err
	}

	job, err := u.jobR.Job(ctx, scripts.JobID(jobID))
	if err != nil {
		return JobDTO{}, err
	}

	if job.OwnerID() != scripts.UserID(userID) {
		return JobDTO{}, scripts.ErrPermissionDenied
	}
	return JobToDTO(*job)
}
