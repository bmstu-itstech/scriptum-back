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
		l.WarnContext(ctx, "user not found")
		return nil, fmt.Errorf("%w: %s", ports.ErrUserNotFound, string(id))
	}
	if err != nil {
		l.ErrorContext(ctx, "failed to select user", slog.String("error", err.Error()))
		return nil, err
	}

	return userRowToDomain(rU)
}

func (r *Repository) Users(ctx context.Context) ([]*entity.User, error) {
	l := r.l.With(slog.String("op", "postgres.Repository.Users"))
	rUs, err := r.selectUserRows(ctx, r.db)
	if err != nil {
		l.ErrorContext(ctx, "failed to select users", slog.String("error", err.Error()))
		return nil, err
	}
	return userRowsToDomain(rUs)
}

func (r *Repository) UserByEmail(ctx context.Context, email string) (*entity.User, error) {
	l := r.l.With(
		slog.String("op", "postgres.Repository.UserByEmail"),
		slog.String("email", email),
	)

	rU, err := r.selectUserRowByEmail(ctx, r.db, email)
	if errors.Is(err, sql.ErrNoRows) {
		l.WarnContext(ctx, "user not found")
		return nil, fmt.Errorf("%w: %s", ports.ErrUserNotFound, email)
	}
	if err != nil {
		l.ErrorContext(ctx, "failed to select user", slog.String("error", err.Error()))
		return nil, err
	}

	return userRowToDomain(rU)
}
