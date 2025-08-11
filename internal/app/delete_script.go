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
		u.logger.Error("failed to delete script", "err1", err)
		return err
	}

	file, err := u.fileR.File(ctx, script.FileID())
	if err != nil {
		u.logger.Error("failed to delete script", "err2", err)
		return err
	}

	if !script.IsAvailableFor(scripts.UserID(actorID)) {
		u.logger.Error("failed to delete script", "err3", err)
		return scripts.ErrPermissionDenied
	}

	err = u.scriptR.Delete(ctx, scripts.ScriptID(scriptID))
	if err != nil {
		u.logger.Error("failed to delete script", "err4", err)
		return err
	}

	err = u.fileR.Delete(ctx, scripts.ScriptID(script.FileID()))
	if err != nil {
		u.logger.Error("failed to delete script", "err5", err)
		_, err := u.scriptR.Create(ctx, &script.ScriptPrototype)
		if err != nil {
			u.logger.Error("failed to restore script", "err6", err)
			return err
		}
		return err
	}

	err = u.manager.Delete(ctx, file.URL())
	if err != nil {
		u.logger.Error("failed to delete script", "err7", err)
		url := file.URL()
		restoredFileId, err := u.fileR.Create(ctx, &url)
		if err != nil {
			u.logger.Error("failed to restore script", "err8", err)
			return err
		}
		newFile, _ := scripts.NewFile(restoredFileId, url)
		newScriptProto, _ := scripts.NewScriptPrototype(
			script.ScriptPrototype.OwnerID(),
			script.ScriptPrototype.Name(),
			script.ScriptPrototype.Desc(),
			script.ScriptPrototype.Visibility(),
			script.ScriptPrototype.Input(),
			script.ScriptPrototype.Output(),
			*newFile,
		)
		_, err = u.scriptR.Create(ctx, newScriptProto)
		if err != nil {
			u.logger.Error("failed to restore script", "err9", err)
			return err
		}
		return err
	}

	return err
}
