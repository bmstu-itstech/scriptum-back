package query

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/response"
	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type GetBlueprintsHandler struct {
	bp ports.BlueprintProvider
	l  *slog.Logger
}

func NewGetBlueprintsHandler(bp ports.BlueprintProvider, l *slog.Logger) GetBlueprintsHandler {
	return GetBlueprintsHandler{bp, l}
}

func (h GetBlueprintsHandler) Handle(ctx context.Context, req request.GetBlueprints) (response.GetBlueprints, error) {
	l := h.l.With(
		slog.String("op", "app.GetBlueprints"),
		slog.String("uid", req.ActorID),
	)

	l.DebugContext(ctx, "querying blueprints")
	bs, err := h.bp.BlueprintsWithUsers(ctx, value.UserID(req.ActorID))
	if err != nil {
		l.ErrorContext(ctx, "failed to query blueprints", slog.String("error", err.Error()))
		return nil, err
	}
	l.InfoContext(ctx, "got blueprints", slog.Int("count", len(bs)))

	return bs, nil
}
