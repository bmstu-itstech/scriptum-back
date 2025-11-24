package ports

import (
	"context"
	"io"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type FileUploader interface {
	Upload(ctx context.Context, name string, reader io.Reader) (value.FileID, error)
}
