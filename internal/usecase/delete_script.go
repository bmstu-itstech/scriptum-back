package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptDeleteUC struct {
	scriptR scripts.ScriptRepository
	userR   scripts.UserRepository
	manager scripts.Manager
}

func NewScriptDeleteUC(scriptR scripts.ScriptRepository, userR scripts.UserRepository, manager scripts.Manager) (*ScriptDeleteUC, error) {
	if scriptR == nil {
		return nil, scripts.ErrInvalidScriptService
	}
	if userR == nil {
		return nil, scripts.ErrInvalidUserService
	}
	if manager == nil {
		return nil, scripts.ErrInvalidManagerService
	}

	return &ScriptDeleteUC{scriptR: scriptR, userR: userR, manager: manager}, nil
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
