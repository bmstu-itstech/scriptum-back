package local

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

func (s *Storage) Read(ctx context.Context, id value.FileID) (io.ReadCloser, error) {
	l := s.l.With(
		slog.String("op", "local.Storage.Read"),
		slog.String("file_id", string(id)),
	)

	dirPath := filepath.Join(s.dir, string(id))

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, ports.ErrFileNotFound
	}

	if len(entries) == 0 {
		return nil, fmt.Errorf("empty file directory: %s", dirPath)
	}

	entry := entries[0]
	path := filepath.Join(dirPath, entry.Name())

	l = l.With(slog.String("path", path))
	l.DebugContext(ctx, "reading file")

	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			l.ErrorContext(ctx, "file not found")
			return nil, ports.ErrFileNotFound
		}
		l.ErrorContext(ctx, "failed to open file", slog.String("error", err.Error()))
		return nil, err
	}
	l.DebugContext(ctx, "successfully open file")

	return file, nil
}
