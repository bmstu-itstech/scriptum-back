package command

import (
	"context"
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
	bp, err := h.br.Blueprint(ctx, value.BlueprintID(req.BlueprintID))
	if err != nil {
		return err
	}

	if bp.OwnerID() != value.UserID(req.ActorID) {
		return domain.ErrPermissionDenied
	}

	err = h.br.DeleteBlueprint(ctx, value.BlueprintID(req.BlueprintID))
	if err != nil {
		return err
	}

	return nil
}
