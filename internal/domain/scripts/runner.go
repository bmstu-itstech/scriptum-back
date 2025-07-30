package scripts

import "context"

type Runner interface {
	Run(context.Context, *Job) (Result, error)
}
