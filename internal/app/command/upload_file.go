package command

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/app/ports"
)

type UploadFileHandler struct {
	u ports.FileUploader
	l *slog.Logger
}

func NewUploadFileHandler(u ports.FileUploader, l *slog.Logger) UploadFileHandler {
	return UploadFileHandler{u, l}
}

func (h UploadFileHandler) Handle(ctx context.Context, req request.UploadFileRequest) (string, error) {
	l := h.l.With(
		slog.String("op", "app.UploadFile"),
		slog.String("name", req.Name),
	)
	l.DebugContext(ctx, "uploading file")
	id, err := h.u.Upload(ctx, req.Name, req.Reader)
	if err != nil {
		l.ErrorContext(ctx, "failed to upload file", slog.String("error", err.Error()))
		return "", err
	}
	l.InfoContext(ctx, "successfully uploaded file", slog.String("id", string(id)))
	return string(id), nil
}
