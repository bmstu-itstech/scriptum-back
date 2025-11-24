package local

import (
	"log/slog"
	"os"
)

const basePathPerms = 0755

type Storage struct {
	basePath string
	l        *slog.Logger
}

func NewStorage(basePath string, l *slog.Logger) (*Storage, error) {
	err := os.MkdirAll(basePath, basePathPerms)
	if err != nil {
		return nil, err
	}
	return &Storage{basePath, l}, nil
}

func MustNewStorage(basePath string, l *slog.Logger) *Storage {
	s, err := NewStorage(basePath, l)
	if err != nil {
		panic(err)
	}
	return s
}
