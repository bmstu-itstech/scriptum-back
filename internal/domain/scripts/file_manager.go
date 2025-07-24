package scripts

import (
	"context"
	"errors"
)

var ErrFileNotFound = errors.New("file not found")

// URL описывает местонахождение файла.
type URL = string

type FileManager interface {
	Save(context.Context, *File) (URL, error)
	Delete(context.Context, URL) error
}
