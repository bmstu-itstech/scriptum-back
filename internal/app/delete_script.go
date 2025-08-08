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
		u.logger.Error("failed to delete script", "err", err)
		return err
	}

	file, err := u.fileR.File(ctx, script.FileID())
	if err != nil {
		u.logger.Error("failed to delete script", "err", err)
		return err
	}

	if !script.IsAvailableFor(scripts.UserID(actorID)) {
		u.logger.Error("failed to delete script", "err", err)
		return scripts.ErrPermissionDenied
	}

	err = u.scriptR.Delete(ctx, scripts.ScriptID(scriptID))
	if err != nil {
		u.logger.Error("failed to delete script", "err", err)
		return err
	}

	err = u.fileR.Delete(ctx, scripts.ScriptID(scriptID))
	if err != nil {
		u.logger.Error("failed to delete script", "err", err)
		_, err := u.scriptR.Create(ctx, &script.ScriptPrototype)
		if err != nil {
			u.logger.Error("failed to restore script", "err", err)
			return err
		}
		return err
	}

	err = u.manager.Delete(ctx, file.URL())
	if err != nil {
		u.logger.Error("failed to delete script", "err", err)
		url := file.URL()
		_, err := u.fileR.Create(ctx, &url)
		if err != nil {
			u.logger.Error("failed to restore script", "err", err)
			return err
		}
		_, err = u.scriptR.Create(ctx, &script.ScriptPrototype)
		if err != nil {
			u.logger.Error("failed to restore script", "err", err)
			return err
		}
		return err
	}

	return err
}
