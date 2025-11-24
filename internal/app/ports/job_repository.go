package ports

import (
	"context"
	"errors"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

var ErrJobAlreadyExists = errors.New("job already exists")

type JobRepository interface {
	SaveJob(ctx context.Context, job *entity.Job) error
	UpdateJob(ctx context.Context, id value.JobID, updateFn func(ctx2 context.Context, job *entity.Job) error) error
}
