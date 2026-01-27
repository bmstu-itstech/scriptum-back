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

type DeleteUserHandler struct {
	ur ports.UserRepository
	l  *slog.Logger
}

func NewDeleteUserHandler(ur ports.UserRepository, l *slog.Logger) DeleteUserHandler {
	return DeleteUserHandler{ur, l}
}

func (h DeleteUserHandler) Handle(ctx context.Context, req request.DeleteUser) error {
	l := h.l.With(
		slog.String("op", "app.DeleteUser"),
		slog.String("actor_id", req.ActorID),
		slog.String("user_id", req.UID),
	)

	actor, err := h.ur.User(ctx, value.UserID(req.ActorID))
	if err != nil {
		l.ErrorContext(ctx, "failed to fetch user", slog.String("error", err.Error()))
		return err
	}
	if actor.Role() != value.RoleAdmin {
		l.WarnContext(ctx, "actor is not admin")
		return domain.ErrPermissionDenied
	}

	err = h.ur.DeleteUser(ctx, value.UserID(req.UID))
	if errors.Is(err, ports.ErrUserNotFound) {
		l.WarnContext(ctx, "user does not exist")
		return err
	} else if err != nil {
		l.ErrorContext(ctx, "failed to delete user", slog.String("error", err.Error()))
		return err
	}

	return nil
}
