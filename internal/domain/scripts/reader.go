package scripts

import "context"

type Reader interface {
	ReadFile(context.Context, string, string, string) (*File, error)
}
