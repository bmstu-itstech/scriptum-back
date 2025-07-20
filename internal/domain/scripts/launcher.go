package scripts

import "context"

type Launcher interface {
	Launch(context.Context, Job, []Field, Email, bool) error
}
