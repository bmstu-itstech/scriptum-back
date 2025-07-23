package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetScriptUC struct {
	scriptR scripts.ScriptRepository
}

func NewGetScript(scriptR scripts.ScriptRepository) (*GetScriptUC, error) {
	if scriptR == nil {
		panic(scripts.ErrInvalidScriptRepository)
	}
	return &GetScriptUC{scriptR: scriptR}, nil
}

func (u *GetScriptUC) Script(ctx context.Context, scriptId int) (ScriptDTO, error) {
	// добавить проверку, что это скрипт конкретного юзера
	s, err := u.scriptR.Script(ctx, scripts.ScriptID(scriptId))
	if err != nil {
		return ScriptDTO{}, err
	}
	return ScriptToDTO(s), nil
}
