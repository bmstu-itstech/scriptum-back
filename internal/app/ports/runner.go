package ports

import (
	"context"
	"io"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type Runner interface {
	Build(ctx context.Context, archive io.Reader, id value.BlueprintID) (value.ImageTag, error)
	Run(ctx context.Context, image value.ImageTag, input []value.Value) (value.Result, error)
	Cleanup(ctx context.Context, image value.ImageTag) error
}
