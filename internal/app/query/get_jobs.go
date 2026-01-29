package query

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/app/dto"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/response"
	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type GetJobsHandler struct {
	jp ports.JobProvider
	l  *slog.Logger
}

func NewGetJobsHandler(jp ports.JobProvider, l *slog.Logger) GetJobsHandler {
	return GetJobsHandler{jp, l}
}

func (h GetJobsHandler) Handle(ctx context.Context, req request.GetJobs) (response.GetJobs, error) {
	l := h.l.With(
		slog.String("op", "app.GetJobs"),
		slog.String("uid", req.ActorID),
	)
	var optState *value.JobState
	if req.State != nil {
		state, err := value.JobStateFromString(*req.State)
		if err != nil {
			return nil, err
		}
		l = l.With("state", state.String())
		optState = &state
	}

	l.DebugContext(ctx, "querying jobs")

	var jobs []dto.Job
	var err error
	if optState == nil {
		jobs, err = h.jp.UserJobs(ctx, value.UserID(req.ActorID))
	} else {
		jobs, err = h.jp.UserJobsWithState(ctx, value.UserID(req.ActorID), *optState)
	}
	if err != nil {
		l.ErrorContext(ctx, "failed to query jobs", slog.String("error", err.Error()))
		return response.GetJobs{}, err
	}
	l.InfoContext(ctx, "got jobs", slog.Int("count", len(jobs)))

	return jobs, nil
}
