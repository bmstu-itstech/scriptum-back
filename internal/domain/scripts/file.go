package scripts

import "fmt"

type File struct {
	name    string // len(name) > 0
	content []byte
}

func (f *File) Name() string {
	return f.name
}

func (f *File) Content() []byte {
	return f.content
}

func NewFile(name string, content []byte) (*File, error) {
	if name == "" {
		return nil, fmt.Errorf("%w: expected not empty filename", ErrInvalidInput)
	}

	return &File{
		name:    name,
		content: content,
	}, nil
}
