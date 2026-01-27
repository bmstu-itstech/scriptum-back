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
	bp ports.BlueprintProvider
	jr ports.JobRepository
	jp ports.JobPublisher
	l  *slog.Logger
}

func NewStartJobHandler(
	bp ports.BlueprintProvider,
	jr ports.JobRepository,
	jp ports.JobPublisher,
	l *slog.Logger,
) StartJobHandler {
	return StartJobHandler{bp, jr, jp, l}
}

func (h StartJobHandler) Handle(ctx context.Context, req request.StartJob) (string, error) {
	l := h.l.With(
		slog.String("op", "app.StartJob"),
		slog.String("blueprint_id", req.BlueprintID),
		slog.String("uid", req.ActorID),
	)
	l.DebugContext(ctx, "starting job", "input", fmt.Sprintf("%+v", req.Values))

	blueprint, err := h.bp.Blueprint(ctx, value.BlueprintID(req.BlueprintID))
	if err != nil {
		l.WarnContext(ctx, "blueprint not found")
		return "", err
	}

	if !blueprint.IsAvailableFor(value.UserID(req.ActorID)) {
		l.WarnContext(ctx, "blueprint is not available")
		return "", domain.ErrPermissionDenied
	}

	in, err := dto.ValuesFromDTOs(req.Values)
	if err != nil {
		l.InfoContext(ctx, "invalid input values", slog.String("error", err.Error()))
		return "", err
	}

	job, err := blueprint.AssembleJob(value.UserID(req.ActorID), in)
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
