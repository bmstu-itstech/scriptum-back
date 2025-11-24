package ports

import (
	"context"
	"errors"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

var ErrJobNotFound = errors.New("job not found")

type JobProvider interface {
	Job(ctx context.Context, id value.JobID) (*entity.Job, error)
	UserJobs(ctx context.Context, uid value.UserID) ([]*entity.Job, error)
	UserJobsWithState(ctx context.Context, uid value.UserID, state value.JobState) ([]*entity.Job, error)
}
