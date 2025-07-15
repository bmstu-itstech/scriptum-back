package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/bmstu-itstech/scriptum-back/internal/service"
)

type ScriptDeleteUC struct {
	scriptS service.ScriptService
}

func NewScriptDeleteUC(scriptS service.ScriptService) (*ScriptDeleteUC, error) {
	if scriptS == nil {
		return nil, scripts.ErrInvalidScriptService
	}

	return &ScriptDeleteUC{scriptS: scriptS}, nil
}

func (u *ScriptDeleteUC) DeleteScript(ctx context.Context, scriptID uint32) error {
	err := u.scriptS.DeleteScript(ctx, scripts.ScriptID(scriptID))
	// логи
	return err
}
