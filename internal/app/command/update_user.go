package command

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/app/dto"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/response"
	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type UpdateUserHandler struct {
	ur ports.UserRepository
	ph ports.PasswordHasher
	l  *slog.Logger
}

func NewUpdateUserHandler(ur ports.UserRepository, ph ports.PasswordHasher, l *slog.Logger) UpdateUserHandler {
	return UpdateUserHandler{ur, ph, l}
}

func (h UpdateUserHandler) Handle(ctx context.Context, req request.UpdateUser) (response.UpdateUser, error) {
	l := h.l.With(
		slog.String("op", "app.UpdateUser"),
		slog.String("actor_id", req.ActorID),
		slog.String("user_id", req.UserID),
	)

	actor, err := h.ur.User(ctx, value.UserID(req.ActorID))
	if err != nil {
		l.ErrorContext(ctx, "failed to fetch user", slog.String("error", err.Error()))
		return response.UpdateUser{}, err
	}
	if actor.Role() != value.RoleAdmin {
		l.WarnContext(ctx, "actor is not admin")
		return response.UpdateUser{}, domain.ErrPermissionDenied
	}

	var ret *entity.User
	err = h.ur.UpdateUser(ctx, value.UserID(req.UserID), func(inner context.Context, u *entity.User) error {
		if pEmail := req.Email; pEmail != nil {
			email, errTx := value.EmailFromString(*pEmail)
			if errTx != nil {
				l.WarnContext(ctx, "failed to validate email", slog.String("error", errTx.Error()))
				return errTx
			}
			errTx = u.SetEmail(email)
			if errTx != nil {
				return errTx
			}
			l.InfoContext(ctx, "updated user email", slog.String("email", email.String()))
		}

		if pPassword := req.Password; pPassword != nil {
			password, errTx := h.ph.Hash(*pPassword)
			if errTx != nil {
				l.WarnContext(ctx, "failed to validate password", slog.String("error", errTx.Error()))
				return errTx
			}
			errTx = u.SetPassword(password)
			if errTx != nil {
				return errTx
			}
			l.InfoContext(ctx, "updated user password")
		}

		if pName := req.Name; pName != nil {
			errTx := u.SetName(*pName)
			if errTx != nil {
				return errTx
			}
			l.InfoContext(ctx, "updated user name", slog.String("name", *pName))
		}

		if pRole := req.Role; pRole != nil {
			role, errTx := value.RoleFromString(*pRole)
			if errTx != nil {
				l.WarnContext(ctx, "failed to validate role", slog.String("error", errTx.Error()))
				return errTx
			}
			if actor.ID() == u.ID() {
				l.WarnContext(ctx, "user can not change self role")
				return domain.ErrPermissionDenied
			}
			errTx = u.SetRole(role)
			if errTx != nil {
				return errTx
			}
			l.InfoContext(ctx, "updated user role", slog.String("role", *pRole))
		}
		ret = u
		return nil
	})
	if err != nil {
		l.ErrorContext(ctx, "failed to update user", slog.String("error", err.Error()))
		return response.UpdateUser{}, err
	}
	return dto.UserToDTO(ret), nil
}
