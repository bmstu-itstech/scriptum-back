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

//nolint:gocognit // Метод состоит из повторяющихся блоков кода.
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
		l.InfoContext(ctx, "actor is not admin")
		return response.UpdateUser{}, domain.ErrPermissionDenied
	}

	var ret *entity.User
	err = h.ur.UpdateUser(ctx, value.UserID(req.UserID), func(inner context.Context, u *entity.User) error {
		if pEmail := req.Email; pEmail != nil {
			errTx := h.updateEmail(inner, l, u, *pEmail)
			if errTx != nil {
				return errTx
			}
		}

		if pPassword := req.Password; pPassword != nil {
			errTx := h.updatePassword(inner, l, u, *pPassword)
			if errTx != nil {
				return errTx
			}
		}

		if pName := req.Name; pName != nil {
			errTx := u.SetName(*pName)
			if errTx != nil {
				l.ErrorContext(ctx, "failed to set user name", slog.String("error", errTx.Error()))
				return errTx
			}
			l.InfoContext(ctx, "updated user name", slog.String("name", *pName))
		}

		if pRole := req.Role; pRole != nil {
			errTx := h.updateRole(ctx, l, u, actor, *pRole)
			if errTx != nil {
				return errTx
			}
		}

		ret = u
		return nil
	})
	if err != nil {
		return response.UpdateUser{}, err
	}
	return dto.UserToDTO(ret), nil
}

func (h UpdateUserHandler) updateEmail(
	ctx context.Context, l *slog.Logger, u *entity.User, emailStr string,
) error {
	email, errTx := value.EmailFromString(emailStr)
	if errTx != nil {
		l.InfoContext(ctx, "failed to validate email", slog.String("error", errTx.Error()))
		return errTx
	}
	errTx = u.SetEmail(email)
	if errTx != nil {
		l.ErrorContext(ctx, "failed to set user email", slog.String("error", errTx.Error()))
		return errTx
	}
	l.InfoContext(ctx, "updated user email", slog.String("email", email.String()))
	return nil
}

func (h UpdateUserHandler) updatePassword(
	ctx context.Context, l *slog.Logger, u *entity.User, passwordStr string,
) error {
	password, errTx := h.ph.Hash(passwordStr)
	if errTx != nil {
		l.InfoContext(ctx, "failed to validate password", slog.String("error", errTx.Error()))
		return errTx
	}
	errTx = u.SetPassword(password)
	if errTx != nil {
		l.ErrorContext(ctx, "failed to set user password", slog.String("error", errTx.Error()))
		return errTx
	}
	l.InfoContext(ctx, "updated user password")
	return nil
}

func (h UpdateUserHandler) updateRole(
	ctx context.Context, l *slog.Logger, u *entity.User, actor *entity.User, roleStr string,
) error {
	role, errTx := value.RoleFromString(roleStr)
	if errTx != nil {
		l.InfoContext(ctx, "failed to validate role", slog.String("error", errTx.Error()))
		return errTx
	}
	if actor.ID() == u.ID() {
		l.InfoContext(ctx, "user can not change self role")
		return domain.ErrPermissionDenied
	}
	errTx = u.SetRole(role)
	if errTx != nil {
		l.ErrorContext(ctx, "failed to set user role", slog.String("error", errTx.Error()))
		return errTx
	}
	l.InfoContext(ctx, "updated user role", slog.String("role", role.String()))
	return nil
}
