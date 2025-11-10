package service

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/google/uuid"
)

type SystemManager struct {
	directory   string
	maxFileSize int64
	l           *slog.Logger
}

func NewSystemManager(dir string, maxFileSize int64, l *slog.Logger) (*SystemManager, error) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	return &SystemManager{
		directory:   dir,
		maxFileSize: maxFileSize,
		l:           l,
	}, nil
}

func (s *SystemManager) Save(_ context.Context, name string, content io.Reader) (scripts.URL, error) {
	s.l.Info("saving file", "name", name)
	dirName := generateDirname(s.directory)

	s.l.Debug("creating dir")
	err := os.MkdirAll(dirName, 0755)
	s.l.Debug("created dir", "err", err.Error())
	if err != nil {
		s.l.Error("failed to create dir", "err", err.Error())
		return "", err
	}

	fileName := fmt.Sprintf("%s/%s", dirName, name)

	s.l.Debug("creating file")
	file, err := os.Create(fileName)
	s.l.Debug("created file", "err", err.Error())
	if err != nil {
		s.l.Error("failed to create file", "err", err.Error())
		return "", err
	}

	s.l.Debug("copying file")
	_, err = io.Copy(file, content)
	s.l.Debug("copied file", "err", err.Error())
	if err != nil {
		s.l.Error("failed to copy file", "err", err.Error())
		return "", err
	}

	s.l.Debug("closing file")
	if err := file.Close(); err != nil {
		s.l.Error("failed to close file", "err", err.Error())
		return "", err
	}

	s.l.Debug("file saved", "fileName", fileName)
	return scripts.URL(fileName), nil
}

func (s *SystemManager) Delete(_ context.Context, path scripts.URL) error {
	s.l.Info("deleting file", "path", path)
	s.l.Debug("deleting file", "path", path)
	err := os.Remove(path)
	s.l.Debug("deleted file", "err", err.Error())
	if err != nil {
		s.l.Error("failed to delete file", "err", err.Error())
		return fmt.Errorf("%w: %s (%w)", scripts.ErrFileNotFound, path, err)
	}
	s.l.Info("file deleted", "path", path)
	return nil
}

func (s *SystemManager) Read(ctx context.Context, path scripts.URL) (scripts.FileData, error) {
	s.l.Info("reading file", "path", path)
	filePath := string(path)

	fileName := filepath.Base(filePath)
	s.l.Debug("reading file", "filePath", filePath)
	file, err := os.Open(filePath)
	s.l.Debug("read file", "err", err.Error())
	if err != nil {
		s.l.Error("failed to read file", "err", err.Error())
		return scripts.FileData{}, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	s.l.Debug("created file data")
	return scripts.FileData{
		Reader: file,
		Name:   fileName,
	}, nil
}

func generateDirname(dir string) string {
	dirName := fmt.Sprintf("%s/%s", dir, uuid.New().String())

	if len(dirName) > scripts.FileURLMaxLen {
		runes := []rune(dirName)
		if len(runes) > scripts.FileURLMaxLen {
			runes = runes[:scripts.FileURLMaxLen]
		}
		dirName = string(runes)
	}

	return dirName
}
