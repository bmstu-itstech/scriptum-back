package scripts

import "context"

type Dispatcher interface {
	Start(context.Context, *Job) error
}
