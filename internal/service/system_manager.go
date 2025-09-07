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

	limitedReader := io.LimitReader(content, s.maxFileSize+1)

	written, err := io.Copy(file, limitedReader)
	if err != nil {
		return "", err
	}

	if written > s.maxFileSize {
		_ = os.RemoveAll(dirName)
		return "", fmt.Errorf("%w:file size exceeds limit of %d bytes", scripts.ErrFileUpload, s.maxFileSize)
	}

	if err := file.Close(); err != nil {
		return "", err
	}

	return scripts.URL(fileName), nil
}

func (s *SystemManager) Delete(_ context.Context, path scripts.URL) error {
	dir := filepath.Dir(path)

	err := os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("%w: %s (%w)", scripts.ErrFileNotFound, path, err)
	}
	return nil
}

func (s *SystemManager) CreateSandbox(ctx context.Context, mainFile scripts.URL, extraFiles []scripts.URL) (scripts.URL, error) {
	dirName := generateDirname(s.directory)

	if err := os.MkdirAll(dirName, 0755); err != nil {
		return "", err
	}

	mainDst := filepath.Join(dirName, filepath.Base(mainFile))
	if err := copyFile(mainFile, mainDst, s.maxFileSize); err != nil {
		_ = os.RemoveAll(dirName)
		return "", err
	}

	for _, f := range extraFiles {
		dst := filepath.Join(dirName, filepath.Base(f))
		if err := copyFile(f, dst, s.maxFileSize); err != nil {
			_ = os.RemoveAll(dirName)
			return "", err
		}
	}

	return mainDst, nil
}

func copyFile(src, dst string, maxSize int64) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer out.Close()

	limited := io.LimitReader(in, maxSize+1)
	written, err := io.Copy(out, limited)
	if err != nil {
		_ = os.Remove(dst)
		return err
	}

	if written > maxSize {
		_ = os.Remove(dst)
		return fmt.Errorf("%w: file size exceeds limit of %d bytes", scripts.ErrFileUpload, maxSize)
	}

	return nil
}

func generateDirname(dir string) string {
	filename := fmt.Sprintf("%s/%s", dir, uuid.New().String())

	if len(filename) > scripts.FileURLMaxLen {
		runes := []rune(filename)
		if len(runes) > scripts.FileURLMaxLen {
			runes = runes[:scripts.FileURLMaxLen]
		}
		filename = string(runes)
	}

	return filename
}
