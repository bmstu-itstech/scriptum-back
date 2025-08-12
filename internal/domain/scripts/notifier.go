package scripts

import "context"

type Notifier interface {
	Notify(context.Context, *Job, Email) error
}
