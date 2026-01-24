package command

import (
	"context"
	"errors"

	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/response"
	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type CreateUserHandler struct {
	ur ports.UserRepository
	ph ports.PasswordHasher
	l  *slog.Logger
}

func NewCreateUserHandler(ur ports.UserRepository, ph ports.PasswordHasher, l *slog.Logger) CreateUserHandler {
	return CreateUserHandler{ur, ph, l}
}

func (h CreateUserHandler) Handle(ctx context.Context, req request.CreateUser) (response.CreateUser, error) {
	l := h.l.With(
		slog.String("op", "app.CreateUser"),
		slog.String("actor_id", req.ActorID),
	)

	actor, err := h.ur.User(ctx, value.UserID(req.ActorID))
	if err != nil {
		l.ErrorContext(ctx, "failed to fetch user", slog.String("error", err.Error()))
		return response.CreateUser{}, err
	}
	if actor.Role() != value.RoleAdmin {
		l.WarnContext(ctx, "actor is not admin")
		return response.CreateUser{}, domain.ErrPermissionDenied
	}

	email, err := value.EmailFromString(req.Email)
	if err != nil {
		l.InfoContext(ctx, "failed to parse email", slog.String("error", err.Error()))
		return response.CreateUser{}, err
	}

	passhash, err := h.ph.Hash(req.Password)
	if err != nil {
		l.ErrorContext(ctx, "failed to hash password", slog.String("error", err.Error()))
		return response.CreateUser{}, err
	}

	role, err := value.RoleFromString(req.Role)
	if err != nil {
		l.InfoContext(ctx, "failed to parse role", slog.String("error", err.Error()))
		return response.CreateUser{}, err
	}

	user, err := entity.NewUser(req.Name, email, passhash, role)
	if err != nil {
		l.ErrorContext(ctx, "failed to create user", slog.String("error", err.Error()))
		return response.CreateUser{}, err
	}

	err = h.ur.SaveUser(ctx, user)
	if errors.Is(err, ports.ErrUserAlreadyExists) {
		l.WarnContext(ctx, "user already exists", slog.String("error", err.Error()))
		return response.CreateUser{}, ports.ErrUserAlreadyExists
	} else if err != nil {
		l.ErrorContext(ctx, "failed to save user", slog.String("error", err.Error()))
		return response.CreateUser{}, err
	}
	l.InfoContext(ctx, "successfully created user", slog.String("id", string(user.ID())))

	return response.CreateUser{
		UID: string(user.ID()),
	}, nil
}
