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

func NewGetScriptUC(scriptR scripts.ScriptRepository, logger *slog.Logger) GetScriptUC {
	return GetScriptUC{
		scriptR: scriptR,
		logger:  logger,
	}
}

func (u *GetScriptUC) Script(ctx context.Context, actorID int64, scriptID int32) (ScriptDTO, error) {
	u.logger.Info("get script", "scriptID", scriptID)
	s, err := u.scriptR.Script(ctx, scripts.ScriptID(scriptID))
	u.logger.Info("got script", "script", s, "err", err)

	if err != nil {
		u.logger.Error("failed to get script", "err", err.Error())
		return ScriptDTO{}, err
	}

	u.logger.Debug("check script availability", "script", s, "actorID", actorID)
	u.logger.Debug("is available", "is", s.IsAvailableFor(scripts.UserID(actorID)))
	if !s.IsAvailableFor(scripts.UserID(actorID)) {
		u.logger.Error("failed to get script", "err", scripts.ErrPermissionDenied.Error())
		return ScriptDTO{}, scripts.ErrPermissionDenied
	}

	u.logger.Info("convert to dto")
	return ScriptToDTO(s)
}
