package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type SearchJobsUC struct {
	jobS  scripts.JobRepository
	userS scripts.UserRepository
}

func NewSearchJobsUC(jobS scripts.JobRepository, userS scripts.UserRepository) (*SearchJobsUC, error) {
	if jobS == nil {
		return nil, scripts.ErrInvalidJobService
	}
	if userS == nil {
		return nil, scripts.ErrInvalidUserService
	}
	return &SearchJobsUC{
		jobS:  jobS,
		userS: userS,
	}, nil
}

func (u *SearchJobsUC) SearchJobs(ctx context.Context, userID uint32, substr string) ([]JobDTO, error) {
	user, err := u.userS.User(ctx, scripts.UserID(userID))
	if err != nil {
		return nil, err
	}
	jobs, err := u.jobS.SearchPublicJobs(ctx, substr)
	if err != nil {
		return nil, err
	}

	if !user.IsAdmin() {
		userJobs, err := u.jobS.SearchUserJobs(ctx, scripts.UserID(userID), substr)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, userJobs...)
	}

	dto := make([]JobDTO, 0, len(jobs))
	for _, job := range jobs {
		dto = append(dto, JobToDTO(job))
	}
	return dto, nil
}
