package command

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/app/dto"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/response"
	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type CreateBlueprintHandler struct {
	br ports.BlueprintRepository
	up ports.UserProvider
	l  *slog.Logger
}

func NewCreateBlueprintHandler(br ports.BlueprintRepository, up ports.UserProvider, l *slog.Logger) CreateBlueprintHandler {
	return CreateBlueprintHandler{br, up, l}
}

func (h CreateBlueprintHandler) Handle(ctx context.Context, req request.CreateBlueprint) (response.CreateBlueprint, error) {
	l := h.l.With(
		slog.String("op", "app.CreateBlueprint"),
		slog.String("uid", req.UID),
	)

	user, err := h.up.User(ctx, value.UserID(req.UID))
	if err != nil {
		l.WarnContext(ctx, "unknown user creating blueprint", slog.String("error", err.Error()))
		return "", err
	}
	l = l.With(slog.String("role", user.Role().String()))

	l.DebugContext(ctx, "creating blueprint", "request", req)

	input, err := dto.FieldsFromDTOs(req.In)
	if err != nil {
		l.InfoContext(ctx, "failed to convert input to dto.Out", slog.String("error", err.Error()))
		return "", err
	}
	output, err := dto.FieldsFromDTOs(req.Out)
	if err != nil {
		l.InfoContext(ctx, "failed to convert output to dto.Out", slog.String("error", err.Error()))
		return "", err
	}

	blueprint, err := entity.NewBlueprint(
		value.UserID(req.UID),
		value.FileID(req.ArchiveID),
		req.Name,
		req.Desc,
		user.BlueprintVisibility(),
		input,
		output,
	)
	if err != nil {
		l.InfoContext(ctx, "failed to create blueprint", slog.String("error", err.Error()))
		return "", err
	}

	err = h.br.SaveBlueprint(ctx, blueprint)
	if err != nil {
		l.ErrorContext(ctx, "failed to save blueprint", slog.String("error", err.Error()))
		return "", err
	}
	l.InfoContext(ctx, "successfully created blueprint", slog.String("id", string(blueprint.ID())))

	return string(blueprint.ID()), nil
}
