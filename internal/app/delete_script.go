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
	script, err := u.scriptR.Script(ctx, scripts.ScriptID(scriptID))
	if err != nil {
		return err
	}

	file, err := u.fileR.File(ctx, script.FileID())
	if err != nil {
		return err
	}

	if !script.IsAvailableFor(scripts.UserID(actorID)) {
		return scripts.ErrPermissionDenied
	}

	err = u.scriptR.Delete(ctx, scripts.ScriptID(scriptID))
	if err != nil {
		return err
	}

	err = u.fileR.Delete(ctx, scripts.ScriptID(scriptID))
	if err != nil {
		return err
	}

	err = u.manager.Delete(ctx, file.URL())
	if err != nil {
		return err
	}

	return err
}
