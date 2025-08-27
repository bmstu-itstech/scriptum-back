package scripts

import (
	"fmt"
	"unicode/utf8"
)

type FileID int64

const FileURLMaxLen = 200

type File struct {
	id  FileID
	url string
}

func (f *File) ID() FileID {
	return f.id
}

func (f *File) URL() string {
	return f.url
}

func NewFile(id FileID, url string) (*File, error) {
	if url == "" {
		return nil, fmt.Errorf("%w: file url must not be empty", ErrInvalidInput)
	}

	if utf8.RuneCountInString(url) > FileURLMaxLen {
		return nil, fmt.Errorf(
			"%w: invalid File: expected len(url) <= %d, got len(url) = %d",
			ErrInvalidInput, FileURLMaxLen, utf8.RuneCountInString(url),
		)
	}

	return &File{
		id:  id,
		url: url,
	}, nil
}
