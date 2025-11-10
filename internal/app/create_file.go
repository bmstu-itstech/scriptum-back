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
	u.logger.Debug("create file debug", "req", req, "ctx", ctx)

	u.logger.Debug("manager save file")
	url, err := u.manager.Save(ctx, req.Name, req.Reader)
	u.logger.Debug("manager saved file", "url", url, "err", err.Error())
	if err != nil {
		u.logger.Error("failed to save file", "err", err.Error())
		return 0, err
	}

	u.logger.Debug("file repository create file")
	fileID, err := u.fileR.Create(ctx, &url)
	u.logger.Debug("file repository created file", "fileID", fileID, "err", err.Error())

	if err != nil {
		u.logger.Error("failed to save (create) file", "err", err.Error())
		u.logger.Debug("manager delete file")
		err_ := u.manager.Delete(ctx, url)
		u.logger.Debug("manager deleted file", "err", err_.Error())
		if err_ != nil {
			u.logger.Error("failed to delete file after error in saving it with manager", "err", err_)
			return 0, err
		}
		return 0, err
	}

	u.logger.Info("file created", "url", url)
	return int32(fileID), nil
}
