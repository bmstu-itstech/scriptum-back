package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zhikh23/pgutils"

	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

func (r *Repository) SaveUser(ctx context.Context, u *entity.User) error {
	rU := userRowFromDomain(u)
	err := r.insertUserRow(ctx, r.db, rU)
	if pgutils.IsUniqueViolationError(err) {
		return fmt.Errorf("%w: %s", ports.ErrUserAlreadyExists, string(u.ID()))
	}
	return err
}

func (r *Repository) UpdateUser(ctx context.Context, uid value.UserID, updateFn func(inner context.Context, u *entity.User) error) error {
	return pgutils.RunTx(ctx, r.db, func(tx *sqlx.Tx) error {
		rU, err := r.selectUserRow(ctx, tx, string(uid))
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("%w: %s", ports.ErrUserNotFound, string(uid))
		} else if err != nil {
			return fmt.Errorf("failed to select for update: %w", err)
		}

		u, err := userRowToDomain(rU)
		if err != nil {
			return err
		}

		err = updateFn(ctx, u)
		if err != nil {
			return err
		}

		rU = userRowFromDomain(u)
		return r.updateUserRow(ctx, tx, rU)
	})
}

func (r *Repository) DeleteUser(ctx context.Context, uid value.UserID) error {
	err := r.softDeleteUserRow(ctx, r.db, string(uid))
	if errors.Is(err, sql.ErrNoRows) {
		return ports.ErrUserNotFound
	}
	return err
}
