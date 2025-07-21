package scripts

import "context"

type JobRepository interface {
	PostJob(context.Context, Job, ScriptID) (JobID, error)
	CloseJob(context.Context, JobID, *Result) error
}
