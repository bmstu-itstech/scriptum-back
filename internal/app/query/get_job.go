package query

import (
	"context"
	"errors"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/app/dto"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/response"
	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type GetJobHandler struct {
	jp ports.JobProvider
	l  *slog.Logger
}

func NewGetJobHandler(jp ports.JobProvider, l *slog.Logger) GetJobHandler {
	return GetJobHandler{jp, l}
}

func (h GetJobHandler) Handle(ctx context.Context, req request.GetJob) (response.GetJob, error) {
	l := h.l.With(
		slog.String("op", "app.GetJob"),
		slog.String("job_id", req.JobID),
		slog.Int64("uid", req.UID),
	)

	l.DebugContext(ctx, "querying job")
	job, err := h.jp.Job(ctx, value.JobID(req.JobID))
	if err != nil {
		if errors.Is(err, ports.ErrJobNotFound) {
			l.WarnContext(ctx, "job not found")
		} else {
			l.ErrorContext(ctx, "failed to query job", slog.String("error", err.Error()))
		}
		return response.GetJob{}, err
	}

	if job.OwnerID() != value.UserID(req.UID) {
		l.WarnContext(ctx, "user does not own job", slog.Int64("owner_id", int64(job.OwnerID())))
		return response.GetJob{}, domain.ErrPermissionDenied
	}
	l.InfoContext(ctx, "got job", slog.String("state", job.State().String()))

	return dto.JobToDTO(job), nil
}
