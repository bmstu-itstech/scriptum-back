package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptDeleteUC struct {
	scriptR scripts.ScriptRepository
	fileR   scripts.FileRepository
	userP   scripts.UserProvider
	manager scripts.FileManager
	logger  *slog.Logger
}

func NewScriptDeleteUC(
	scriptR scripts.ScriptRepository,
	userR scripts.UserProvider,
	manager scripts.FileManager,
	fileR scripts.FileRepository,
	logger *slog.Logger,
) ScriptDeleteUC {
	return ScriptDeleteUC{
		scriptR: scriptR,
		userP:   userR,
		fileR:   fileR,
		manager: manager,
		logger:  logger,
	}
}

func (u *ScriptDeleteUC) DeleteScript(ctx context.Context, actorID uint32, scriptID uint32) error {
	u.logger.Info("deleting script", "actorID", actorID)
	script, err := u.scriptR.Script(ctx, scripts.ScriptID(scriptID))
	if err != nil {
		u.logger.Error("failed to get script", "err", err)
		return err
	}

	file, err := u.fileR.File(ctx, script.MainFileID())
	if err != nil {
		u.logger.Error("failed to get file", "err", err)
		return err
	}

	if !script.IsAvailableFor(scripts.UserID(actorID)) {
		u.logger.Error("script is not available to delete by user", "err", err)
		return scripts.ErrPermissionDenied
	}

	err = u.scriptR.Delete(ctx, scripts.ScriptID(scriptID))
	if err != nil {
		u.logger.Error("failed to delete script", "err", err)
		return err
	}

	err = u.fileR.Delete(ctx, scripts.ScriptID(script.MainFileID()))
	if err != nil {
		u.logger.Error("failed to delete main file while deleting script", "err", err)
		_, err := u.scriptR.Restore(ctx, &script)
		if err != nil {
			u.logger.Error("failed to restore script while deleting script", "err", err)
			return err
		}
		return err
	}

	for _, e := range script.ExtraFileIDs() {
		err := u.fileR.Delete(ctx, scripts.ScriptID(e))
		if err != nil {
			u.logger.Error("failed to delete extra file while deleting script", "err", err)
			_, err := u.scriptR.Restore(ctx, &script)
			if err != nil {
				u.logger.Error("failed to restore script while deleting script", "err", err)
				return err
			}
			return err
		}
	}

	err = u.manager.Delete(ctx, file.URL())
	if err != nil {
		u.logger.Error("failed to delete file from system", "err", err)
		_, err := u.fileR.Restore(ctx, file)
		if err != nil {
			u.logger.Error("failed to restore file while deleting file from system", "err", err)
			return err
		}

		_, err = u.scriptR.Restore(ctx, &script)
		if err != nil {
			u.logger.Error("failed to restore script while deleting file from system", "err9", err)
			return err
		}
		return err
	}

	return err
}
