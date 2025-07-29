package scripts

import "context"

type Runner interface {
	Run(ctx context.Context, job *Job, path URL, expected []Field) (Result, error)
}
