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
	dirName := generateDirname(s.directory)

	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		return "", err
	}

	fileName := fmt.Sprintf("%s/%s", dirName, name)

	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, content)
	if err != nil {
		return "", err
	}

	if err := file.Close(); err != nil {
		return "", err
	}

	return scripts.URL(fileName), nil
}

func (s *SystemManager) Delete(_ context.Context, path scripts.URL) error {
	err := os.Remove(path)
	if err != nil {
		return fmt.Errorf("%w: %s (%w)", scripts.ErrFileNotFound, path, err)
	}
	return nil
}

func (s *SystemManager) Read(ctx context.Context, path scripts.URL) (scripts.FileData, error) {
	filePath := string(path)

	fileName := filepath.Base(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		return scripts.FileData{}, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

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
