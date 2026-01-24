package command

import (
	"context"
	"errors"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/response"
	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain"
)

type LoginHandler struct {
	ur ports.UserProvider
	ph ports.PasswordHasher
	ts ports.TokenService
	l  *slog.Logger
}

func NewLoginHandler(ur ports.UserProvider, ph ports.PasswordHasher, ts ports.TokenService, l *slog.Logger) LoginHandler {
	return LoginHandler{ur, ph, ts, l}
}

func (h *LoginHandler) Handle(ctx context.Context, req request.Login) (response.LoginResponse, error) {
	l := h.l.With(
		slog.String("op", "app.Login"),
		slog.String("email", req.Email),
	)

	user, err := h.ur.UserByEmail(ctx, req.Email)
	if errors.Is(err, ports.ErrUserNotFound) {
		l.WarnContext(ctx, "user not found")
		return response.LoginResponse{}, domain.ErrInvalidCredentials
	} else if err != nil {
		l.ErrorContext(ctx, "failed to get user", slog.String("error", err.Error()))
		return response.LoginResponse{}, err
	}

	if !h.ph.Verify(req.Password, user.PasswordHash()) {
		l.WarnContext(ctx, "password mismatch")
		return response.LoginResponse{}, domain.ErrInvalidCredentials
	}

	token, err := h.ts.GenerateToken(ctx, user.ID())
	if err != nil {
		l.ErrorContext(ctx, "failed to generate token", slog.String("error", err.Error()))
		return response.LoginResponse{}, err
	}

	return response.LoginResponse{AccessToken: string(token)}, nil
}
