package service

import "github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"

type FileUploader interface {
	scripts.Uploader
}

type FileUpload struct{}

func NewFileUploader() *FileUpload {
	return &FileUpload{}
}

func (f *FileUpload) upload(file scripts.File) (scripts.Path, error) {

	return "", nil
}
