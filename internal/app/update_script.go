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
	script, err := u.scriptR.Script(ctx, scripts.ScriptID(req.ID))
	if err != nil {
		return err
	}

	if !script.IsAvailableFor(scripts.UserID(actorID)) {
		return scripts.ErrPermissionDenied
	}

	proto, err := DTOToScript(req)
	if err != nil {
		return err
	}

	err = u.scriptR.Update(ctx, proto)
	return err
}
