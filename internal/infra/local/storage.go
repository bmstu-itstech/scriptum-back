package local

import (
	"errors"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/config"
)

const basePathPerms = 0755

type Storage struct {
	dir string
	l   *slog.Logger
}

func NewStorage(cfg config.Storage, l *slog.Logger) (*Storage, error) {
	if l == nil {
		return nil, errors.New("nil logger")
	}
	if cfg.BasePath == "" {
		return nil, errors.New("Storage.BasePath is required")
	}
	return &Storage{cfg.BasePath, l}, nil
}

func MustNewStorage(cfg config.Storage, l *slog.Logger) *Storage {
	s, err := NewStorage(cfg, l)
	if err != nil {
		panic(err)
	}
	return s
}
