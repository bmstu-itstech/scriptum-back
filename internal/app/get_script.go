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
	if scriptR == nil {
		panic(scripts.ErrInvalidScriptRepository)
	}
	if logger == nil {
		panic(scripts.ErrInvalidLogger)
	}
	return GetScriptUC{scriptR: scriptR, logger: logger}
}

func (u *GetScriptUC) Script(ctx context.Context, scriptId int) (ScriptDTO, error) {
	s, err := u.scriptR.Script(ctx, scripts.ScriptID(scriptId))
	if err != nil {
		return ScriptDTO{}, err
	}
	return ScriptToDTO(s), nil
}
