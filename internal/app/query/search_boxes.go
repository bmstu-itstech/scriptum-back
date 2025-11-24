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

type SearchBoxesHandler struct {
	bp ports.BoxProvider
	l  *slog.Logger
}

func NewSearchBoxesHandler(bp ports.BoxProvider, l *slog.Logger) SearchBoxesHandler {
	return SearchBoxesHandler{bp, l}
}

func (h SearchBoxesHandler) Handle(ctx context.Context, req request.SearchBoxes) (response.SearchBoxes, error) {
	l := h.l.With(
		slog.String("op", "app.SearchBoxes"),
		slog.Int64("uid", req.UID),
	)

	l.DebugContext(ctx, "querying boxes")
	boxes, err := h.bp.SearchBoxes(ctx, value.UserID(req.UID), req.Name)
	if err != nil {
		l.ErrorContext(ctx, "failed to search boxes", slog.String("error", err.Error()))
		return nil, err
	}
	l.InfoContext(ctx, "found boxes", slog.Int("count", len(boxes)))

	return dto.BoxesToDTOs(boxes), nil
}
