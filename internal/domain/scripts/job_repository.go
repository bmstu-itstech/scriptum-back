package scripts

import "context"

type JobRepository interface {
	// Create(context.Context) (*Job, error)
	// Get(context.Context, JobID) (Job, error)
	Store(context.Context, Job) (JobID, error)
}
