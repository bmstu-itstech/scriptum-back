package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptDeleteUC struct {
	scriptR scripts.ScriptRepository
	userR   scripts.UserRepository
	logger  *slog.Logger
	manager scripts.Manager
}

func NewScriptDeleteUC(scriptR scripts.ScriptRepository, userR scripts.UserRepository, logger *slog.Logger, manager scripts.Manager) ScriptDeleteUC {
	if scriptR == nil {
		panic(scripts.ErrInvalidScriptRepository)
	}
	if userR == nil {
		panic(scripts.ErrInvalidUserRepository)
	}
	if manager == nil {
		panic(scripts.ErrInvalidManagerService)
	}
	if logger == nil {
		panic(scripts.ErrInvalidLogger)
	}

	return ScriptDeleteUC{scriptR: scriptR, userR: userR, logger: logger, manager: manager}
}

func (u *ScriptDeleteUC) DeleteScript(ctx context.Context, actorID uint32, scriptID uint32) error {
	var err error
	user, err := u.userR.User(ctx, scripts.UserID(actorID))
	if err != nil {
		return err
	}
	script, err := u.scriptR.Script(ctx, scripts.ScriptID(scriptID))
	if err != nil {
		return err
	}

	if adm := user.IsAdmin(); adm && script.Visibility() == scripts.VisibilityGlobal || !adm && script.Owner() == actorID {
		err = u.scriptR.DeleteScript(ctx, scripts.ScriptID(scriptID))
		if err != nil {
			return err
		}
		err = u.manager.Delete(ctx, script.Path())
	} else {
		err = scripts.ErrNoAccessToDelete
	}

	// логи
	return err
}
