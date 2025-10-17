package scripts

import (
	"context"
	"io"
)

type FileData struct {
	Reader io.Reader
	Name   string
}

type Launcher interface {
	CreateSandbox(ctx context.Context, mainReader FileData, extraReaders []FileData) (URL, error)
	Run(context.Context, *Job) (Result, error)
	DeleteSandbox(ctx context.Context, path URL) error
}
