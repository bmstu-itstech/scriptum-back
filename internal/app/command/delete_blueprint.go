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
		slog.String("actor_id", req.ActorID),
		slog.String("blueprint_id", req.BlueprintID),
	)

	bp, err := h.br.Blueprint(ctx, value.BlueprintID(req.BlueprintID))
	if errors.Is(err, ports.ErrBlueprintNotFound) {
		l.WarnContext(ctx, "blueprint not found")
		return nil
	} else if err != nil {
		l.ErrorContext(ctx, "could not find blueprint", slog.String("error", err.Error()))
		return err
	}

	if bp.OwnerID() != value.UserID(req.ActorID) {
		l.WarnContext(ctx, "not authorized to delete this blueprint", slog.String("owner_id", string(bp.OwnerID())))
		return domain.ErrPermissionDenied
	}

	err = h.br.DeleteBlueprint(ctx, value.BlueprintID(req.BlueprintID))
	if err != nil {
		l.ErrorContext(ctx, "could not delete blueprint", slog.String("error", err.Error()))
		return err
	}
	l.InfoContext(ctx, "successfully deleted blueprint")

	return nil
}
