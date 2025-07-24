package scripts

import (
	"context"
	"errors"
)

var ErrJobNotFound = errors.New("job not found")

type JobRepository interface {
	Create(context.Context, *JobPrototype) (*Job, error)
	Update(context.Context, *Job) error
	Delete(context.Context, JobID) error

	Job(context.Context, JobID) (*Job, error)
	UserJobs(context.Context, UserID) ([]Job, error)
	UserJobsWithState(context.Context, UserID, JobState) ([]Job, error)
}
