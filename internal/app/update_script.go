package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptUpdateUC struct {
	scriptR scripts.ScriptRepository
	userR   scripts.UserRepository
	logger  *slog.Logger
}

func NewScriptUpdateUC(scriptR scripts.ScriptRepository, userR scripts.UserRepository, logger *slog.Logger) ScriptUpdateUC {
	if scriptR == nil {
		panic(scripts.ErrInvalidScriptRepository)
	}
	if userR == nil {
		panic(scripts.ErrInvalidUserRepository)
	}
	if logger == nil {
		panic(scripts.ErrInvalidLogger)
	}
	return ScriptUpdateUC{scriptR: scriptR, userR: userR, logger: logger}
}

func (u *ScriptUpdateUC) Update(ctx context.Context, actorID uint32, input ScriptDTO) error {
	user, err := u.userR.User(ctx, scripts.UserID(actorID))
	if err != nil {
		return err
	}
	script, err := DTOToScript(input)
	if err != nil {
		return err
	}

	if adm := user.IsAdmin(); adm && script.Visibility() == scripts.VisibilityGlobal || !adm && input.Owner == actorID {
		err = u.scriptR.UpdateScript(ctx, script)
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
