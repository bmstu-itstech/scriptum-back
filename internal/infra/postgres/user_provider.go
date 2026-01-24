package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

func (r *Repository) User(ctx context.Context, id value.UserID) (*entity.User, error) {
	rU, err := r.selectUserRow(ctx, r.db, string(id))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: %s", ports.ErrUserNotFound, string(id))
	}
	if err != nil {
		return nil, err
	}
	return userRowToDomain(rU)
}

func (r *Repository) Users(ctx context.Context) ([]*entity.User, error) {
	rUs, err := r.selectUserRows(ctx, r.db)
	if err != nil {
		return nil, err
	}
	return userRowsToDomain(rUs)
}

func (r *Repository) UserByEmail(ctx context.Context, email string) (*entity.User, error) {
	rU, err := r.selectUserRowByEmail(ctx, r.db, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: %s", ports.ErrUserNotFound, email)
	}
	if err != nil {
		return nil, err
	}
	return userRowToDomain(rU)
}
