package scripts

import "context"

type Notifier interface {
	Notify(context.Context, Result) error
}
