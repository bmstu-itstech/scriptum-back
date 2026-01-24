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

type DeleteBlueprintHandler struct {
	br ports.BlueprintRepository
	l  *slog.Logger
}

func NewDeleteBlueprintHandler(br ports.BlueprintRepository, l *slog.Logger) DeleteBlueprintHandler {
	return DeleteBlueprintHandler{br, l}
}

func (h DeleteBlueprintHandler) Handle(ctx context.Context, req request.DeleteBlueprint) error {
	l := h.l.With(
		slog.String("op", "app.DeleteBlueprint"),
		slog.String("uid", req.ActorID),
	)

	bp, err := h.br.Blueprint(ctx, value.BlueprintID(req.BlueprintID))
	if err != nil {
		if errors.Is(err, ports.ErrBlueprintNotFound) {
			l.WarnContext(ctx, "blueprint not found")
		} else {
			l.ErrorContext(ctx, "failed to query blueprint", slog.String("error", err.Error()))
		}
		return err
	}

	if bp.OwnerID() != value.UserID(req.ActorID) {
		l.WarnContext(ctx, "user can't delete blueprint", slog.String("owner_id", string(bp.OwnerID())))
		return domain.ErrPermissionDenied
	}

	err = h.br.DeleteBlueprint(ctx, value.BlueprintID(req.BlueprintID))
	if err != nil {
		l.ErrorContext(ctx, "failed to delete blueprint", slog.String("error", err.Error()))
		return err
	}
	l.InfoContext(ctx, "blueprint deleted")

	return nil
}
