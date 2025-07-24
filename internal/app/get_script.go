package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetScriptUC struct {
	scriptR scripts.ScriptRepository
	logger  *slog.Logger
}

func NewGetScript(scriptR scripts.ScriptRepository, logger *slog.Logger) GetScriptUC {
	return GetScriptUC{
		scriptR: scriptR,
		logger:  logger,
	}
}

func (u *GetScriptUC) Script(ctx context.Context, actorId int64, scriptId int32) (ScriptDTO, error) {
	s, err := u.scriptR.Script(ctx, scripts.ScriptID(scriptId))
	if err != nil {
		return ScriptDTO{}, err
	}

	if !s.IsAvailableFor(scripts.UserID(actorId)) {
		return ScriptDTO{}, scripts.ErrPermissionDenied
	}

	return ScriptToDTO(s), nil
}
