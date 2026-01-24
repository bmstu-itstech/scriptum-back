package query

import (
	"context"
	"errors"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/app/dto"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/response"
	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type GetBlueprintHandler struct {
	bp ports.BlueprintProvider
	l  *slog.Logger
}

func NewGetBlueprintHandler(bp ports.BlueprintProvider, l *slog.Logger) GetBlueprintHandler {
	return GetBlueprintHandler{bp, l}
}

func (h GetBlueprintHandler) Handle(ctx context.Context, req request.GetBlueprint) (response.GetBlueprint, error) {
	l := h.l.With(
		slog.String("op", "app.GetBlueprint"),
		slog.String("blueprint_id", req.BlueprintID),
		slog.String("uid", req.UID),
	)

	l.DebugContext(ctx, "querying blueprint")
	blueprint, err := h.bp.Blueprint(ctx, value.BlueprintID(req.BlueprintID))
	if err != nil {
		if errors.Is(err, ports.ErrBlueprintNotFound) {
			l.WarnContext(ctx, "blueprint not found")
		} else {
			l.ErrorContext(ctx, "failed to query blueprint", slog.String("error", err.Error()))
		}
		return response.GetBlueprint{}, err
	}

	if !blueprint.IsAvailableFor(value.UserID(req.BlueprintID)) {
		l.WarnContext(ctx, "user can't see the blueprint", slog.String("owner_id", string(blueprint.OwnerID())))
		return response.GetBlueprint{}, domain.ErrPermissionDenied
	}
	l.InfoContext(ctx, "got blueprint")

	return dto.BlueprintToDTO(blueprint), nil
}
