package service

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type FileManager struct {
	directory string
}

func NewFileUploader() (*FileManager, error) {
	return &FileManager{
		directory: "scripts",
	}, nil
}

func (f *FileManager) Upload(_ context.Context, file scripts.File) (scripts.Path, error) {
	dir := filepath.Join(f.directory, file.FileType())

	err := os.MkdirAll(dir, 0755)

	if err != nil {
		return "", err
	}

	mask := rand.Int()
	filename := fmt.Sprintf("%d_%d_%s", mask, time.Now().UnixNano(), filepath.Base(file.Name()))
	filename = filename[:100]

	path := filepath.Join(dir, filename)

	err = os.WriteFile(filename, file.Content(), 0644)
	return path, err
}

func (f *FileManager) Delete(_ context.Context, path scripts.Path) error {
	return os.Remove(path)
}
