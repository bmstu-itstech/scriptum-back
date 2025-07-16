package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptDeleteUC struct {
	scriptS   scripts.ScriptRepository
	userS     scripts.UserRepository
	uploaderS scripts.Uploader
}

func NewScriptDeleteUC(scriptS scripts.ScriptRepository, userS scripts.UserRepository, uploaderS scripts.Uploader) (*ScriptDeleteUC, error) {
	if scriptS == nil {
		return nil, scripts.ErrInvalidScriptService
	}
	if userS == nil {
		return nil, scripts.ErrInvalidUserService
	}
	if uploaderS == nil {
		return nil, scripts.ErrInvalidUploaderService
	}

	return &ScriptDeleteUC{scriptS: scriptS, userS: userS, uploaderS: uploaderS}, nil
}

func (u *ScriptDeleteUC) DeleteScript(ctx context.Context, actorID uint32, scriptID uint32) error {
	var err error
	user, err := u.userS.User(ctx, scripts.UserID(actorID))
	if err != nil {
		return err
	}
	script, err := u.scriptS.Script(ctx, scripts.ScriptID(scriptID))
	if err != nil {
		return err
	}

	if adm := user.IsAdmin(); adm && script.Visibility() == scripts.VisibilityGlobal || !adm && script.Owner() == actorID {
		err = u.scriptS.DeleteScript(ctx, scripts.ScriptID(scriptID))
		if err != nil {
			return err
		}
		err = u.uploaderS.Delete(ctx, script.Path())
	} else {
		err = scripts.ErrNoAccessToDelete
	}

	// логи
	return err
}
