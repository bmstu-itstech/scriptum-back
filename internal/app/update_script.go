package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptUpdateUC struct {
	scriptR scripts.ScriptRepository
	logger  *slog.Logger
}

func NewScriptUpdateUC(
	scriptR scripts.ScriptRepository,
	logger *slog.Logger,
) ScriptUpdateUC {
	return ScriptUpdateUC{
		scriptR: scriptR,
		logger:  logger,
	}
}

func (u *ScriptUpdateUC) UpdateScript(ctx context.Context, actorID int64, req ScriptDTO) error {
	u.logger.Info("updating script ", "req", req)
	script, err := u.scriptR.Script(ctx, scripts.ScriptID(req.ID))
	if err != nil {
		u.logger.Error("failed to update script", "err", err)
		return err
	}

	if !script.IsAvailableFor(scripts.UserID(actorID)) {
		u.logger.Error("failed to update script", "err", scripts.ErrPermissionDenied)
		return scripts.ErrPermissionDenied
	}

	proto, err := DTOToScript(req)
	if err != nil {
		u.logger.Error("failed to update script", "err", err)
		return err
	}

	return u.scriptR.Update(ctx, proto)
}
