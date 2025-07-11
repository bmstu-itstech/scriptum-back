package service

import "github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"

type FileUpload struct{}

func NewFileUploader() (*FileUpload, error) {
	return &FileUpload{}, nil
}

func (f *FileUpload) Upload(file scripts.File) (scripts.Path, error) {

	return "", nil
}
