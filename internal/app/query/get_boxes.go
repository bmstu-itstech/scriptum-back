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

type GetBoxesHandler struct {
	bp ports.BoxProvider
	l  *slog.Logger
}

func NewGetBoxesHandler(bp ports.BoxProvider, l *slog.Logger) GetBoxesHandler {
	return GetBoxesHandler{bp, l}
}

func (h GetBoxesHandler) Handle(ctx context.Context, req request.GetBoxes) (response.GetBoxes, error) {
	l := h.l.With(
		slog.String("op", "app.GetBoxes"),
		slog.Int64("uid", req.UID),
	)

	l.DebugContext(ctx, "querying boxes")
	boxes, err := h.bp.Boxes(ctx, value.UserID(req.UID))
	if err != nil {
		l.ErrorContext(ctx, "failed to query boxes", slog.String("error", err.Error()))
		return nil, err
	}
	l.InfoContext(ctx, "got boxes", slog.Int("count", len(boxes)))

	return dto.BoxesToDTOs(boxes), nil
}
