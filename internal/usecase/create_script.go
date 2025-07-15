package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptCreateUC struct {
	scriptS scripts.ScriptRepository
}

func NewScriptCreateUC(scriptS scripts.ScriptRepository) (*ScriptCreateUC, error) {
	if scriptS == nil {
		return nil, scripts.ErrInvalidScriptService
	}
	return &ScriptCreateUC{scriptS: scriptS}, nil
}

func (u *ScriptCreateUC) CreateScript(ctx context.Context, input ScriptDTO) (uint32, error) {
	script, err := DTOToScript(input)
	if err != nil {
		// логируем ошибку
		return 0, err
	}
	scriptId, err := u.scriptS.CreateScript(ctx, script)
	if err != nil {
		// логируем ошибку
		return 0, err
	}
	return uint32(scriptId), nil
}
