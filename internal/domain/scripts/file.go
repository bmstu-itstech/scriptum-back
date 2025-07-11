package scripts

import (
	"context"
	"mime/multipart"
)

type File struct {
	name     string
	fileType string
	content  []byte
}

type FileReader interface {
	ReadFile(context.Context, multipart.File, *multipart.FileHeader) (*File, error)
}

func (f *File) Name() string {
	return f.name
}

func (f *File) FileType() string {
	return f.fileType
}

func (f *File) Content() []byte {
	return f.content
}

func NewFile(name, fileType string, content []byte) (*File, error) {
	if name == "" {
		return nil, ErrFileNameEmpty
	}
	if fileType == "" {
		return nil, ErrFileTypeEmpty
	}
	if len(content) == 0 {
		return nil, ErrFileContentEmpty
	}

	return &File{
		name:     name,
		fileType: fileType,
		content:  content,
	}, nil
}
