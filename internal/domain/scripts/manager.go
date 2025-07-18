package scripts

import "context"

type Manager interface {
	Upload(context.Context, File) (Path, error)
	Delete(context.Context, Path) error
}
