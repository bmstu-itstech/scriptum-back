package scripts

import "context"

type JobRepository interface {
	Post(context.Context, Job, ScriptID) (JobID, error)
	Update(context.Context, JobID, *Result) error
}
