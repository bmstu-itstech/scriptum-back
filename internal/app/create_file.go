package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type FileCreateUC struct {
	fileR   scripts.FileRepository
	userP   scripts.UserProvider
	manager scripts.FileManager
	logger  *slog.Logger
}

func NewFileCreateUC(
	userP scripts.UserProvider,
	fileR scripts.FileRepository,
	manager scripts.FileManager,
	logger *slog.Logger,
) FileCreateUC {
	return FileCreateUC{
		userP:   userP,
		fileR:   fileR,
		manager: manager,
		logger:  logger,
	}
}

func (u *FileCreateUC) CreateFile(ctx context.Context, req FileDTO) (int32, error) {
	u.logger.Info("create file", "req", req)

	url, err := u.manager.Save(ctx, req.Name, req.Reader)
	if err != nil {
		u.logger.Error("failed to save file", "err", err)
		return 0, err
	}

	fileID, err := u.fileR.Create(ctx, &url)
	if err != nil {
		u.logger.Error("failed to save (create) file", "err", err)
		err := u.manager.Delete(ctx, url)
		if err != nil {
			u.logger.Error("failed to delete file after error in saving it with manager", "err", err)
			return 0, err
		}
		return 0, err
	}

	u.logger.Info("file created", "url", url)
	return int32(fileID), nil
}
