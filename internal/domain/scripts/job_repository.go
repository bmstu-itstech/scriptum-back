package scripts

import "context"

type JobRepository interface {
	StoreJob(context.Context, Job) (JobID, error)
	JobsByScriptID(context.Context, ScriptID) ([]Job, error)
	SearchPublicJobs(context.Context, string) ([]Job, error)
	SearchUserJobs(context.Context, UserID, string) ([]Job, error)
}
