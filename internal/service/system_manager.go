package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/google/uuid"
)

type SystemManager struct {
	directory   string
	maxFileSize int64
}

func NewSystemManager(dir string, maxFileSize int64) (*SystemManager, error) {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	return &SystemManager{
		directory:   dir,
		maxFileSize: maxFileSize,
	}, nil
}

func (s *SystemManager) Save(_ context.Context, name string, content io.Reader) (scripts.URL, error) {
	filename := generateFilename(s.directory, name)

	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}

	limitedReader := io.LimitReader(content, s.maxFileSize+1)

	written, err := io.Copy(file, limitedReader)
	if err != nil {
		return "", err
	}

	if written > s.maxFileSize {
		_ = os.Remove(filename)
		return "", fmt.Errorf("%w:file size exceeds limit of %d bytes", scripts.ErrFileUpload, s.maxFileSize)
	}
	
	if err := file.Close(); err != nil {
		return "", err
	}

	return scripts.URL(filename), nil
}

func (s *SystemManager) Delete(_ context.Context, path scripts.URL) error {
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("%w: %s (%w)", scripts.ErrFileNotFound, path, err)
	}
	return nil
}

func generateFilename(dir, originalName string) string {
	base := filepath.Base(originalName)
	filename := fmt.Sprintf("%s/%s_%s", dir, uuid.New().String(), base)

	if len(filename) > scripts.FileURLMaxLen {
		runes := []rune(filename)
		if len(runes) > scripts.FileURLMaxLen {
			runes = runes[:scripts.FileURLMaxLen]
		}
		filename = string(runes)
	}

	return filename
}
