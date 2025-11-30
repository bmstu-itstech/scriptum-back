package local

import (
	"log/slog"
)

const basePathPerms = 0755

type Storage struct {
	dir string
	l   *slog.Logger
}

func NewStorage(dir string, l *slog.Logger) *Storage {
	return &Storage{dir, l}
}
