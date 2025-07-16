package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptUpdateUC struct {
	scriptS scripts.ScriptRepository
	userS   scripts.UserRepository
}

func NewScriptUpdateUC(scriptS scripts.ScriptRepository, userS scripts.UserRepository) (*ScriptUpdateUC, error) {
	if scriptS == nil {
		return nil, scripts.ErrInvalidScriptService
	}
	if userS == nil {
		return nil, scripts.ErrInvalidUserService
	}
	return &ScriptUpdateUC{scriptS: scriptS, userS: userS}, nil
}

func (u *ScriptUpdateUC) Update(ctx context.Context, actorID uint32, input ScriptDTO) error {
	user, err := u.userS.User(ctx, scripts.UserID(actorID))
	if err != nil {
		return err
	}
	script, err := DTOToScript(input)
	if err != nil {
		return err
	}

	if adm := user.IsAdmin(); adm && script.Visibility() == scripts.VisibilityGlobal || !adm && input.Owner == actorID {
		err = u.scriptS.UpdateScript(ctx, script)
	} else {
		err = scripts.ErrNoAccessToUpdate
	}

	return err
	// логика в том, что по переданному айдишнику
	// будут вставлены
	// новые данные из этой же структуры
	// могут поменяться поля, владелец,
	// путь (возможно видимость)
}
