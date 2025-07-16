package scripts

import "context"

type JobRepository interface {
	Store(context.Context, Job) (JobID, error)
	PublicJobs(context.Context) ([]Job, error)
	UserJobs(context.Context, UserID) ([]Job, error)
	SearchPublicJobs(context.Context, string) ([]Job, error)
	SearchUserJobs(context.Context, UserID, string) ([]Job, error)
}
