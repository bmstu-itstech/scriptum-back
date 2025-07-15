package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/bmstu-itstech/scriptum-back/internal/service"
)

type ScriptRunUC struct {
	scriptS service.ScriptService
}

func NewScriptRunUC(scriptS service.ScriptService) (*ScriptRunUC, error) {
	if scriptS == nil {
		return nil, scripts.ErrInvalidScriptService
	}
	return &ScriptRunUC{scriptS: scriptS}, nil
}

func (s *ScriptRunUC) RunScript(ctx context.Context, scriptID scripts.ScriptID) (scripts.Result, error) {
	
}