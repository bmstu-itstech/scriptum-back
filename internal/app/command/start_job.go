package command

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/app/dto"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type StartJobHandler struct {
	bp ports.BoxProvider
	jr ports.JobRepository
	jp ports.JobPublisher
	l  *slog.Logger
}

func NewStartJobHandler(
	bp ports.BoxProvider,
	jr ports.JobRepository,
	jp ports.JobPublisher,
	l *slog.Logger,
) StartJobHandler {
	return StartJobHandler{bp, jr, jp, l}
}

func (h StartJobHandler) Handle(ctx context.Context, req request.StartJob) (string, error) {
	l := h.l.With(
		slog.String("op", "app.StartJob"),
		slog.String("box_id", req.BoxID),
		slog.Int64("uid", req.UID),
	)
	l.DebugContext(ctx, "starting job", "input", fmt.Sprintf("%+v", req.Values))

	box, err := h.bp.Box(ctx, value.BoxID(req.BoxID))
	if err != nil {
		l.InfoContext(ctx, "box not found")
		return "", err
	}

	if !box.IsAvailableFor(value.UserID(req.UID)) {
		l.WarnContext(ctx, "box is not available")
		return "", domain.ErrPermissionDenied
	}

	in, err := dto.ValuesFromDTOs(req.Values)
	if err != nil {
		l.InfoContext(ctx, "invalid input values", slog.String("error", err.Error()))
		return "", err
	}

	job, err := box.AssembleJob(value.UserID(req.UID), in)
	if err != nil {
		l.InfoContext(ctx, "failed to assemble job", slog.String("error", err.Error()))
		return "", err
	}

	err = h.jr.SaveJob(ctx, job)
	if err != nil {
		l.ErrorContext(ctx, "failed to save job", slog.String("error", err.Error()))
		return "", err
	}

	err = h.jp.PublishJob(ctx, job)
	if err != nil {
		l.ErrorContext(ctx, "failed to publish job", slog.String("error", err.Error()))
		return "", err
	}
	l.InfoContext(ctx, "job published successfully", slog.String("id", string(job.ID())))

	return string(job.ID()), nil
}
