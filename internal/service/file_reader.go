package service

import (
	"io"
	"mime/multipart"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

func ReadFile(file multipart.File, header *multipart.FileHeader) (scripts.File, error) {
	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	content, err := io.ReadAll(file)
	if err != nil {
		return scripts.File{}, err
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return scripts.File{
		Name:    header.Filename,
		Type:    mimeType,
		Content: content,
	}, nil
}
