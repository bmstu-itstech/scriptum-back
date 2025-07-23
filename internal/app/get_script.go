package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetScriptUC struct {
	scriptR scripts.ScriptRepository
	userR   scripts.UserRepository
	logger  *slog.Logger
}

func NewGetScript(
	scriptR scripts.ScriptRepository,
	userR scripts.UserRepository,
	logger *slog.Logger,
) GetScriptUC {
	if scriptR == nil {
		panic(scripts.ErrInvalidScriptRepository)
	}
	if userR == nil {
		panic(scripts.ErrInvalidUserRepository)
	}
	if logger == nil {
		panic(scripts.ErrInvalidLogger)
	}
	return GetScriptUC{scriptR: scriptR, userR: userR, logger: logger}
}

func (u *GetScriptUC) Script(ctx context.Context, userID, scriptId int) (ScriptDTO, error) {
	user, err := u.userR.User(ctx, scripts.UserID(userID))
	if err != nil {
		return ScriptDTO{}, err
	}
	s, err := u.scriptR.Script(ctx, scripts.ScriptID(scriptId))
	if err != nil {
		return ScriptDTO{}, err
	}

	if user.IsAdmin() && s.Visibility() == scripts.VisibilityPrivate ||
		!user.IsAdmin() && s.Owner() != scripts.UserID(userID) {
		return ScriptDTO{}, fmt.Errorf("cannot get someone else's script")
	}

	return ScriptToDTO(s), nil
}
