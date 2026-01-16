package api

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	apiv2 "github.com/bmstu-itstech/scriptum-back/gen/go/api/v2"
	"github.com/bmstu-itstech/scriptum-back/internal/app"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
)

type fileService struct {
	apiv2.UnimplementedFileServiceServer

	app *app.App
	l   *slog.Logger
}

func RegisterFileService(s *grpc.Server, app *app.App, l *slog.Logger) {
	apiv2.RegisterFileServiceServer(s, &fileService{app: app, l: l})
}

func (s fileService) Upload(
	stream grpc.ClientStreamingServer[apiv2.FileUploadRequest, apiv2.FileUploadResponse],
) error {
	l := s.l.With(slog.String("op", "api.Upload"))

	name := ""
	size := uint32(0)
	var buf bytes.Buffer

	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			l.Error("error receiving file", slog.String("error", err.Error()))
			return err
		}

		switch body := req.GetBody().(type) {
		case *apiv2.FileUploadRequest_Meta:
			name = body.Meta.GetName()
		case *apiv2.FileUploadRequest_Chunk:
			chunk := req.GetChunk()
			size += uint32(len(chunk)) //nolint:gosec // len всегда возвращает неотрицательное число
			l.Debug("received chunk", slog.Int("size", len(chunk)))
			if _, err2 := buf.Write(chunk); err2 != nil {
				return status.Error(codes.Internal, err2.Error())
			}
		}
	}

	fileID, err := s.app.Commands.UploadFile.Handle(context.Background(), request.UploadFileRequest{
		Name:   name,
		Reader: &buf,
	})
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	return stream.SendAndClose(&apiv2.FileUploadResponse{
		FileId: fileID,
		Size:   size,
	})
}
