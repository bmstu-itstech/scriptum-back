package scripts

import "context"

type JobRepository interface {
	GetResult(context.Context, JobID) (Result, error)
	GetResultsForUser(context.Context, UserID) ([]Result, error)
	PostJob(context.Context, Job, ScriptID) (JobID, error)
	CloseJob(context.Context, JobID, *Result) error
	JobsByScriptID(context.Context, ScriptID) ([]Job, error)
	SearchJobs(context.Context, UserID, string) ([]Result, error)
}
