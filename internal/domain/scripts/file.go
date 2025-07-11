package scripts

import (
	"io"
	"mime/multipart"
)

type File struct {
	Name    string
	Type    string
	Content []byte
}

func readFile(file multipart.File, header *multipart.FileHeader) (File, error) {
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return File{}, err
	}

	mimeType := header.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return File{
		Name:    header.Filename,
		Type:    mimeType,
		Content: content,
	}, nil
}
