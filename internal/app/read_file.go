package app

import (
	"mime/multipart"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/bmstu-itstech/scriptum-back/internal/service"

	"context"
)

type ReadFileHTTP struct {
	fileReader service.HTTPFileReader
}

func NewReadFileHTTPUseCase(fileReader service.HTTPFileReader) (*ReadFileHTTP, error) {
	return &ReadFileHTTP{fileReader: fileReader}, nil
}

func (rf *ReadFileHTTP) ReadFile(ctx context.Context, file multipart.File, header *multipart.FileHeader) (*scripts.File, error) {
	return rf.fileReader.ReadFile(ctx, file, header)
}
