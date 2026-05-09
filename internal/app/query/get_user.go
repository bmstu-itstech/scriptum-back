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

type GetUserHandler struct {
	up ports.UserProvider
	l  *slog.Logger
}

func NewGetUserHandler(up ports.UserProvider, l *slog.Logger) GetUserHandler {
	return GetUserHandler{up, l}
}

func (h GetUserHandler) Handle(ctx context.Context, req request.GetUser) (response.GetUser, error) {
	l := h.l.With(
		slog.String("op", "app.GetUser"),
		slog.String("user_id", req.UserID),
		slog.String("uid", req.ActorID),
	)

	actor, err := h.up.User(ctx, value.UserID(req.ActorID))
	if errors.Is(err, ports.ErrUserNotFound) {
		l.InfoContext(ctx, "user not found")
		return response.GetUser{}, domain.ErrPermissionDenied
	}
	if err != nil {
		l.InfoContext(ctx, "failed to query user", slog.String("error", err.Error()))
		return response.GetUser{}, err
	}

	if !actor.CanSee(value.UserID(req.UserID)) {
		l.InfoContext(ctx, "user does not have permission to see this user")
		return response.GetUser{}, domain.ErrPermissionDenied
	}

	user, err := h.up.User(ctx, value.UserID(req.UserID))
	if err != nil {
		l.InfoContext(ctx, "failed to query user", slog.String("error", err.Error()))
		return response.GetUser{}, err
	}

	return dto.UserToDTO(user), nil
}
