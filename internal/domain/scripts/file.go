package scripts

import (
	"context"
	"fmt"
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
		return nil, fmt.Errorf("name: expected not empty string, got empty string  %w", ErrFileInvalid)
	}
	if fileType == "" {
		return nil, fmt.Errorf("fileType: expected not empty string, got empty string  %w", ErrFileInvalid)
	}
	if len(content) == 0 {
		return nil, fmt.Errorf("content: expected byte array with at least one elemet, got empty array  %w", ErrFileInvalid)
	}

	return &File{
		name:     name,
		fileType: fileType,
		content:  content,
	}, nil
}
