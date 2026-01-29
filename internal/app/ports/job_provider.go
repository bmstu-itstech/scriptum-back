package ports

import (
	"context"
	"errors"

	"github.com/bmstu-itstech/scriptum-back/internal/app/dto"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

var ErrJobNotFound = errors.New("job not found")

type JobProvider interface {
	Job(ctx context.Context, id value.JobID) (dto.Job, error)
	UserJobs(ctx context.Context, uid value.UserID) ([]dto.Job, error)
	UserJobsWithState(ctx context.Context, uid value.UserID, state value.JobState) ([]dto.Job, error)
}
