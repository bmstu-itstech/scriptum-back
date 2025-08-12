package scripts

import (
	"context"
)

type FileRepository interface {
	Create(context.Context, *URL) (FileID, error)
	Delete(context.Context, ScriptID) error
	Restore(context.Context, *File) (FileID, error)
	File(context.Context, FileID) (*File, error)
}
