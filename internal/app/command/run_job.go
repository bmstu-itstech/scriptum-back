package command

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type RunJobHandler struct {
	r  ports.Runner
	jr ports.JobRepository
	fr ports.FileReader
	l  *slog.Logger
}

func NewRunJobHandler(r ports.Runner, jr ports.JobRepository, fr ports.FileReader, l *slog.Logger) RunJobHandler {
	return RunJobHandler{r, jr, fr, l}
}

func (h RunJobHandler) Handle(ctx context.Context, job request.RunJob) error {
	l := h.l.With(
		slog.String("op", "app.RunJob"),
		slog.String("job_id", job.JobID),
	)
	l.Debug("running job")

	// Две раздельные операции обновления Job так как необходимо достичь состояния Running.

	err := h.jr.UpdateJob(ctx, value.JobID(job.JobID), func(ctx2 context.Context, job *entity.Job) error {
		return job.Run()
	})
	if err != nil {
		l.Error("failed to update job", slog.String("error", err.Error()))
		return err
	}

	var res value.Result
	err = h.jr.UpdateJob(ctx, value.JobID(job.JobID), func(ctx2 context.Context, job *entity.Job) error {
		buildCtx, err2 := h.fr.Read(ctx2, job.ArchiveID())
		if err2 != nil {
			return fmt.Errorf("failed to read build context: %w", err2)
		}
		image, err2 := h.r.Build(ctx2, buildCtx, job.BoxID())
		if err2 != nil {
			return err2
		}
		res, err2 = h.r.Run(ctx2, image, job.Input())
		if err2 != nil {
			return err2
		}
		return job.Finish(res)
	})
	if err != nil {
		l.Error("failed to update job", slog.String("error", err.Error()))
		return err
	}
	l.Info("job successfully run", slog.Int("code", int(res.Code())), slog.String("output", res.Output()))
	return nil
}
