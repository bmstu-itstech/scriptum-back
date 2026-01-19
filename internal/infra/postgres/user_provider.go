package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

func (r *Repository) User(ctx context.Context, id value.UserID) (*entity.User, error) {
	l := r.l.With(
		slog.String("op", "postgres.Repository.User"),
		slog.String("id", string(id)),
	)

	rU, err := r.selectUserRow(ctx, r.db, string(id))
	if errors.Is(err, sql.ErrNoRows) {
		l.WarnContext(ctx, "user not found", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %s", ports.ErrUserNotFound, string(id))
	}
	if err != nil {
		l.ErrorContext(ctx, "failed to select user", slog.String("error", err.Error()))
		return nil, err
	}

	return userRowToDomain(rU)
}

func (r *Repository) UserByEmail(ctx context.Context, email value.Email) (*entity.User, error) {
	l := r.l.With(
		slog.String("op", "postgres.Repository.UserByEmail"),
		slog.String("email", email.String()),
	)

	rU, err := r.selectUserRowByEmail(ctx, r.db, email.String())
	if errors.Is(err, sql.ErrNoRows) {
		l.WarnContext(ctx, "user not found", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%w: %s", ports.ErrUserNotFound, email.String())
	}
	if err != nil {
		l.ErrorContext(ctx, "failed to select user", slog.String("error", err.Error()))
		return nil, err
	}

	return userRowToDomain(rU)
}
