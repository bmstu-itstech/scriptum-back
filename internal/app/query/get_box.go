package query

import (
	"context"
	"errors"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/app/dto"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/response"
	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type GetBoxHandler struct {
	bp ports.BoxProvider
	l  *slog.Logger
}

func NewGetBoxHandler(bp ports.BoxProvider, l *slog.Logger) GetBoxHandler {
	return GetBoxHandler{bp, l}
}

func (h GetBoxHandler) Handle(ctx context.Context, req request.GetBox) (response.GetBox, error) {
	l := h.l.With(
		slog.String("op", "app.GetBox"),
		slog.String("box_id", req.BoxID),
		slog.Int64("uid", req.UID),
	)

	l.DebugContext(ctx, "querying box")
	box, err := h.bp.Box(ctx, value.BoxID(req.BoxID))
	if err != nil {
		if errors.Is(err, ports.ErrBoxNotFound) {
			l.WarnContext(ctx, "box not found")
		} else {
			l.ErrorContext(ctx, "failed to query box", slog.String("error", err.Error()))
		}
		return response.GetBox{}, err
	}
	l.InfoContext(ctx, "got box")

	return dto.BoxToDTO(box), nil
}
