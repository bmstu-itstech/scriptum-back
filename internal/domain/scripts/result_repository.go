package scripts

import "context"

type ResultRepository interface {
	GetResult(context.Context, JobID) (Result, error)
	UserResults(context.Context, UserID) ([]Result, error)
	SearchResult(context.Context, UserID, string) ([]Result, error)
}
