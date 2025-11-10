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
	u.logger.Debug("getting script", "scriptID", scriptID)
	script, err := u.scriptR.Script(ctx, scripts.ScriptID(scriptID))
	u.logger.Debug("got script", "script", script, "err", err.Error())
	if err != nil {
		u.logger.Error("failed to get script", "err", err.Error())
		return err
	}

	u.logger.Debug("getting file", "fileID", script.MainFileID())
	file, err := u.fileR.File(ctx, script.MainFileID())
	u.logger.Debug("got file", "file", file, "err", err.Error())
	if err != nil {
		u.logger.Error("failed to get file", "err", err.Error())
		return err
	}

	u.logger.Debug("script is available to delete by user", "actorID", actorID, "is", script.IsAvailableFor(scripts.UserID(actorID)))
	if !script.IsAvailableFor(scripts.UserID(actorID)) {
		u.logger.Error("script is not available to delete by user", "err", scripts.ErrPermissionDenied.Error())
		return scripts.ErrPermissionDenied
	}

	u.logger.Debug("deleting script", "scriptID", scriptID)
	err = u.scriptR.Delete(ctx, scripts.ScriptID(scriptID))
	u.logger.Debug("deleted script", "err", err.Error())
	if err != nil {
		u.logger.Error("failed to delete script", "err", err.Error())
		return err
	}

	u.logger.Debug("deleting main file", "fileID", script.MainFileID())
	err = u.fileR.Delete(ctx, scripts.ScriptID(script.MainFileID()))
	u.logger.Debug("deleted main file", "err", err.Error())
	if err != nil {
		u.logger.Error("failed to delete main file while deleting script", "err", err.Error())
		u.logger.Debug("restoring script", "script", script)
		_, err := u.scriptR.Restore(ctx, &script)
		u.logger.Debug("restored script", "err", err.Error())
		if err != nil {
			u.logger.Error("failed to restore script while deleting script", "err", err.Error())
			return err
		}
		return err
	}

	u.logger.Debug("deleting extra files", "fileIDs", script.ExtraFileIDs())
	for _, e := range script.ExtraFileIDs() {
		u.logger.Debug("deleting extra file", "fileID", e)
		err := u.fileR.Delete(ctx, scripts.ScriptID(e))
		u.logger.Debug("deleted extra file", "err", err.Error())
		if err != nil {
			u.logger.Error("failed to delete extra file while deleting script", "err", err.Error())
			u.logger.Debug("restoring script", "script", script)
			_, err := u.scriptR.Restore(ctx, &script)
			u.logger.Debug("restored script", "err", err.Error())
			if err != nil {
				u.logger.Error("failed to restore script while deleting script", "err", err.Error())
				return err
			}
			return err
		}
	}

	u.logger.Debug("deleting file from system", "file", file)
	err = u.manager.Delete(ctx, file.URL())
	u.logger.Debug("deleted file from system", "err", err.Error())
	if err != nil {
		u.logger.Error("failed to delete file from system", "err", err.Error())
		u.logger.Debug("restoring file", "file", file)
		_, err := u.fileR.Restore(ctx, file)
		u.logger.Debug("restored file", "err", err.Error())
		if err != nil {
			u.logger.Error("failed to restore file while deleting file from system", "err", err.Error())
			return err
		}

		u.logger.Debug("restoring script", "script", script)
		_, err = u.scriptR.Restore(ctx, &script)
		u.logger.Debug("restored script", "err", err.Error())
		if err != nil {
			u.logger.Error("failed to restore script while deleting file from system", "err", err.Error())
			return err
		}
		return err
	}

	u.logger.Info("deleted file from system", "file", file)
	return err
}
