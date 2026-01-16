package postgres

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/jmoiron/sqlx"

	"github.com/bmstu-itstech/scriptum-back/internal/config"
)

type Repository struct {
	db *sqlx.DB
	l  *slog.Logger
}

func NewRepository(cfg config.Postgres, l *slog.Logger) (*Repository, error) {
	if l == nil {
		return nil, errors.New("nil logger")
	}

	db, err := sqlx.Connect("postgres", cfg.URI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}
	return &Repository{db, l}, nil
}

func MustNewRepository(cfg config.Postgres, l *slog.Logger) *Repository {
	r, err := NewRepository(cfg, l)
	if err != nil {
		panic(err)
	}
	return r
}
