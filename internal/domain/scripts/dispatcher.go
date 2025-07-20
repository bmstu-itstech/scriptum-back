package scripts

import "context"

type Dispatcher interface {
	Launch(context.Context, LaunchRequest) error
}
