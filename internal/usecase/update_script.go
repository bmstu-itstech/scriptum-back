package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptUpdateUC struct {
	scriptS scripts.ScriptRepository
}

func NewScriptUpdateUC(scriptS scripts.ScriptRepository) (*ScriptUpdateUC, error) {
	if scriptS == nil {
		return nil, scripts.ErrInvalidScriptService
	}
	return &ScriptUpdateUC{scriptS: scriptS}, nil
}

func (u *ScriptUpdateUC) Update(ctx context.Context, input ScriptDTO) error {
	script, err := DTOToScript(input)
	if err != nil {
		return err
	}
	// логика в том, что по переданному айдишнику 
	// будут вставлены
	// новые данные из этой же структуры
	// могут поменяться поля, владелец, 
	// путь (возможно видимость)
	return u.scriptS.UpdateScript(ctx, script)
}
