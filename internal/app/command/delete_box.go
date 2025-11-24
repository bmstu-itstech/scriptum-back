package command

import (
	"context"
	"errors"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type DeleteBoxHandler struct {
	br ports.BoxRepository
	l  *slog.Logger
}

func NewDeleteBoxHandler(br ports.BoxRepository, l *slog.Logger) DeleteBoxHandler {
	return DeleteBoxHandler{br, l}
}

func (h DeleteBoxHandler) Handle(ctx context.Context, req request.DeleteBox) error {
	l := h.l.With(
		slog.String("op", "app.DeleteBox"),
		slog.Int64("uid", req.UID),
	)

	box, err := h.br.Box(ctx, value.BoxID(req.BoxID))
	if err != nil {
		if errors.Is(err, ports.ErrBoxNotFound) {
			l.Warn("box not found")
		} else {
			l.Error("failed to query box", slog.String("error", err.Error()))
		}
		return err
	}

	if box.OwnerID() != value.UserID(req.UID) {
		l.WarnContext(ctx, "user can't delete box", slog.Int64("owner_id", int64(box.OwnerID())))
		return domain.ErrPermissionDenied
	}

	err = h.br.DeleteBox(ctx, value.BoxID(req.BoxID))
	if err != nil {
		l.ErrorContext(ctx, "failed to delete box", slog.String("error", err.Error()))
		return err
	}
	l.InfoContext(ctx, "box deleted")

	return nil
}
