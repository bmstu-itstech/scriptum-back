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

type SearchBlueprintsHandler struct {
	bp ports.BlueprintProvider
	l  *slog.Logger
}

func NewSearchBlueprintsHandler(bp ports.BlueprintProvider, l *slog.Logger) SearchBlueprintsHandler {
	return SearchBlueprintsHandler{bp, l}
}

func (h SearchBlueprintsHandler) Handle(ctx context.Context, req request.SearchBlueprints) (response.SearchBlueprints, error) {
	l := h.l.With(
		slog.String("op", "app.SearchBlueprints"),
		slog.String("uid", req.UID),
	)

	l.DebugContext(ctx, "querying blueprints")
	bs, err := h.bp.SearchBlueprints(ctx, value.UserID(req.UID), req.Name)
	if err != nil {
		l.ErrorContext(ctx, "failed to search blueprints", slog.String("error", err.Error()))
		return nil, err
	}
	l.InfoContext(ctx, "found blueprints", slog.Int("count", len(bs)))

	return dto.BlueprintsToDTOs(bs), nil
}
