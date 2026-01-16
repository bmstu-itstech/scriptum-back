package local

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

func (s *Storage) Upload(ctx context.Context, name string, reader io.Reader) (value.FileID, error) {
	l := s.l.With(
		slog.String("op", "local.Storage.Upload"),
		slog.String("name", name),
	)

	fileID := value.NewFileID()
	l = l.With(slog.String("file_id", string(fileID)))

	dirPath := filepath.Join(s.dir, string(fileID))

	err := os.MkdirAll(dirPath, basePathPerms)
	if err != nil {
		l.ErrorContext(ctx, "failed to create directory", slog.String("error", err.Error()))
		return "", err
	}

	filePath := filepath.Join(dirPath, name)

	l = l.With(slog.String("path", filePath))
	l.DebugContext(ctx, "saving file")

	file, err := os.Create(filePath)
	if err != nil {
		l.ErrorContext(ctx, "failed to create file", slog.String("error", err.Error()))
		return "", err
	}
	defer func() { _ = file.Close() }()

	_, err = io.Copy(file, reader)
	if err != nil {
		l.ErrorContext(ctx, "failed to copy file", slog.String("error", err.Error()))
		err2 := os.Remove(dirPath)
		if err2 != nil {
			l.ErrorContext(ctx, "failed to remove file", slog.String("error", err.Error()))
		}
		return "", err
	}
	l.DebugContext(ctx, "successfully uploaded file")

	return fileID, nil
}
