package scripts

import "context"

type Launcher interface {
	Launch(context.Context, Job) (Result, error)
}
