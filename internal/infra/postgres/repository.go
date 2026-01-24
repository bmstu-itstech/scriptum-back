package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/bmstu-itstech/scriptum-back/internal/config"
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(cfg config.Postgres) (*Repository, error) {
	db, err := sqlx.Connect("postgres", cfg.URI)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}
	return &Repository{db}, nil
}

func MustNewRepository(cfg config.Postgres) *Repository {
	r, err := NewRepository(cfg)
	if err != nil {
		panic(err)
	}
	return r
}
