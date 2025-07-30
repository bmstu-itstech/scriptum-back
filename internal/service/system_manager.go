package service

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/google/uuid"
)

type SystemManager struct {
	directory string
}

func NewSystemManager(dir string) (*SystemManager, error) {
	err := os.MkdirAll(dir, 0755)

	if err != nil {
		return nil, err
	}

	return &SystemManager{
		directory: dir,
	}, nil
}

func (s *SystemManager) Save(_ context.Context, file *scripts.File) (scripts.URL, error) {
	filename := generateFilename(s.directory, file.Name())

	err := os.WriteFile(filename, file.Content(), 0644)
	return scripts.URL(filename), err
}

func (s *SystemManager) Delete(_ context.Context, path scripts.URL) error {
	return os.Remove(path)
}

func generateFilename(dir, originalName string) string {
	base := filepath.Base(originalName)
	filename := fmt.Sprintf("%s/%s_%s", dir, uuid.New().String(), base)

	if len(filename) > scripts.ScriptURLMaxLen {
		runes := []rune(filename)
		if len(runes) > scripts.ScriptURLMaxLen {
			runes = runes[:scripts.ScriptURLMaxLen]
		}
		filename = string(runes)
	}

	return filename
}
