package scripts

import "context"

type ResultRepository interface {
	JobResult(context.Context, JobID) (Result, error)
	UserResults(context.Context, UserID) ([]Result, error)
	SearchResult(context.Context, UserID, string) ([]Result, error)
}
