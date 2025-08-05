package scripts

import (
	"context"
)

type FileRepository interface {
	Create(context.Context, *URL) (FileID, error)
	Delete(context.Context, ScriptID) error

	File(context.Context, FileID) (*File, error)
}
