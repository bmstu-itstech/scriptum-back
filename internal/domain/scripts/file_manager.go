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
	Save(context.Context, string, io.Reader) (URL, error)
	Delete(context.Context, URL) error
}
