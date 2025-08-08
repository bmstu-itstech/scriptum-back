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

func (u *GetScriptUC) Script(ctx context.Context, actorId int64, scriptID int32) (ScriptDTO, error) {
	u.logger.Info("get script", "scriptID", scriptID)
	s, err := u.scriptR.Script(ctx, scripts.ScriptID(scriptID))
	if err != nil {
		u.logger.Error("failed to get script", "err", err)
		return ScriptDTO{}, err
	}

	if !s.IsAvailableFor(scripts.UserID(actorId)) {
		u.logger.Error("failed to get script", "err", scripts.ErrPermissionDenied)
		return ScriptDTO{}, scripts.ErrPermissionDenied
	}

	return ScriptToDTO(s)
}
