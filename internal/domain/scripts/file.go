package scripts

import (
	"fmt"
)

type FileID int64

const FileURLMaxLen = 300

type File struct {
	id     FileID
	url    string
	isMain bool
}

func (f *File) ID() FileID {
	return f.id
}

func (f *File) URL() string {
	return f.url
}

func (f *File) IsMain() bool {
	return f.isMain
}

func NewFile(id FileID, url string, isMain bool) (*File, error) {
	if url == "" {
		return nil, fmt.Errorf("%w: file url must not be empty", ErrInvalidInput)
	}

	if len(url) > FileURLMaxLen {
		return nil, fmt.Errorf(
			"%w: invalid File: expected len(url) <= %d, got len(url) = %d",
			ErrInvalidInput, FileURLMaxLen, len(url),
		)
	}

	return &File{
		id:     id,
		url:    url,
		isMain: isMain,
	}, nil
}
