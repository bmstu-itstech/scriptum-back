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

type GetUsersHandler struct {
	up ports.UserProvider
	l  *slog.Logger
}

func NewGetUsersHandler(up ports.UserProvider, l *slog.Logger) GetUsersHandler {
	return GetUsersHandler{up, l}
}

func (h GetUsersHandler) Handle(ctx context.Context, req request.GetUsers) (response.GetUsers, error) {
	actor, err := h.up.User(ctx, value.UserID(req.ActorID))
	if errors.Is(err, ports.ErrUserNotFound) {
		return nil, domain.ErrPermissionDenied
	}
	if err != nil {
		return nil, err
	}

	if actor.Role() != value.RoleAdmin {
		return nil, domain.ErrPermissionDenied
	}

	users, err := h.up.Users(ctx)
	if err != nil {
		return nil, err
	}

	return dto.UsersToDTOs(users), nil
}
