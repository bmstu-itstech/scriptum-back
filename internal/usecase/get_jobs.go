package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetJobsUC struct {
	jobS  scripts.JobRepository
	userS scripts.UserRepository
}

func NewGetJobsUC(jobS scripts.JobRepository, userS scripts.UserRepository) (*GetJobsUC, error) {
	if jobS == nil {
		return nil, scripts.ErrInvalidJobService
	}
	if userS == nil {
		return nil, scripts.ErrInvalidUserService
	}
	return &GetJobsUC{jobS: jobS, userS: userS}, nil
}

func (u *GetJobsUC) GetJobs(ctx context.Context, userID uint32) ([]JobDTO, error) {
	jobs, err := u.jobS.PublicJobs(ctx)
	if err != nil {
		return nil, err
	}

	user, err := u.userS.User(ctx, scripts.UserID(userID))
	if err != nil {
		return nil, err
	}

	if !user.IsAdmin() {
		userJobs, err := u.jobS.UserJobs(ctx, scripts.UserID(userID))
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, userJobs...)
	}

	dto := make([]JobDTO, 0, len(jobs))
	for _, v := range jobs {
		dto = append(dto, JobToDTO(v))
	}

	return dto, nil
}
