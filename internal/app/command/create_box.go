package command

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/app/dto"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type CreateBoxHandler struct {
	br  ports.BoxRepository
	iap ports.IsAdminProvider
	l   *slog.Logger
}

func NewCreateBoxHandler(br ports.BoxRepository, iap ports.IsAdminProvider, l *slog.Logger) CreateBoxHandler {
	return CreateBoxHandler{br, iap, l}
}

func (h CreateBoxHandler) Handle(ctx context.Context, req request.CreateBox) (string, error) {
	l := h.l.With(
		slog.String("op", "app.CreateBox"),
		slog.Int64("uid", req.UID),
	)

	isAdmin, err := h.iap.IsAdmin(ctx, value.UserID(req.UID))
	if err != nil {
		l.ErrorContext(ctx, "failed to check author isAdmin", slog.String("error", err.Error()))
		return "", err
	}
	l = l.With(slog.Bool("is_admin", isAdmin))

	l.DebugContext(ctx, "creating box", "request", req)

	var vis value.Visibility
	if isAdmin {
		vis = value.VisibilityPublic
	} else {
		vis = value.VisibilityPrivate
	}

	input, err := dto.FieldsFromDTOs(req.Input)
	if err != nil {
		l.InfoContext(ctx, "failed to convert input to dto.Out", slog.String("error", err.Error()))
		return "", err
	}
	output, err := dto.FieldsFromDTOs(req.Output)
	if err != nil {
		l.InfoContext(ctx, "failed to convert output to dto.Out", slog.String("error", err.Error()))
		return "", err
	}

	box, err := entity.NewBox(
		value.UserID(req.UID),
		value.FileID(req.ArchiveID),
		req.Name,
		req.Desc,
		vis,
		input,
		output,
	)
	if err != nil {
		l.InfoContext(ctx, "failed to create box", slog.String("error", err.Error()))
		return "", err
	}

	err = h.br.SaveBox(ctx, box)
	if err != nil {
		l.ErrorContext(ctx, "failed to save box", slog.String("error", err.Error()))
		return "", err
	}
	l.InfoContext(ctx, "successfully created box", slog.String("id", string(box.ID())))

	return string(box.ID()), nil
}
