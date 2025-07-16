package scripts

import "context"

type Uploader interface {
	Upload(context.Context, File) (Path, error)
	Delete(context.Context, Path) error
}
