package dto

import "io"

type File struct {
	Name   string
	Reader io.Reader
}
