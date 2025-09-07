package scripts

import (
	"context"
	"errors"
	"io"
)

var ErrFileNotFound = errors.New("file not found")

var ErrFileUpload = errors.New("cannot upload file")

// URL описывает местонахождение файла.
type URL = string

type FileManager interface {
	Save(ctx context.Context, name string, reader io.Reader) (URL, error)
	Delete(ctx context.Context, url URL) error
	CreateSandbox(ctx context.Context, mainFile File, extraFiles []File) (URL, error)
}
