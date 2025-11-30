package tests

import (
	"context"
	"errors"
	"io"

	apiv2 "github.com/bmstu-itstech/scriptum-back/gen/go/api/v2"
)

const chunkSize int = 4 * 1024 // 4 Kb

func UploadFile(
	ctx context.Context,
	fs apiv2.FileServiceClient,
	filename string,
	data io.Reader,
) (string, error) {
	stream, err := fs.Upload(ctx)
	if err != nil {
		return "", err
	}

	err = stream.Send(&apiv2.FileUploadRequest{
		Body: &apiv2.FileUploadRequest_Meta{
			Meta: &apiv2.FileMeta{
				Name: filename,
			},
		},
	})
	if err != nil {
		return "", err
	}

	var n int
	for buf := make([]byte, chunkSize); err == nil; n, err = data.Read(buf) {
		err = stream.Send(&apiv2.FileUploadRequest{
			Body: &apiv2.FileUploadRequest_Chunk{
				Chunk: buf[:n],
			},
		})
		if err != nil {
			return "", err
		}
	}
	// err != nil
	if !errors.Is(err, io.EOF) {
		return "", err
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return "", err
	}

	return res.FileId, nil
}
