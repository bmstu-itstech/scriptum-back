package ports

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
)

type JobPublisher interface {
	PublishJob(ctx context.Context, job *entity.Job) error
}
