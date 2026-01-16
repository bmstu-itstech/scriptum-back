package ports

import (
	"context"
	"errors"
	"io"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

var ErrFileNotFound = errors.New("file not found")

type FileReader interface {
	Read(ctx context.Context, id value.FileID) (io.ReadCloser, error)
}
