package service

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type FileUpload struct{}

func NewFileUploader() (*FileUpload, error) {
	return &FileUpload{}, nil
}

func (f *FileUpload) Upload(_ context.Context, file scripts.File) (scripts.Path, error) {

	return "", nil
}
