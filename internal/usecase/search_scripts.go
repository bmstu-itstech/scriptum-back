package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptSearchUC struct {
	scriptS scripts.ScriptRepository
}

func NewScriptSearchUC(scriptS scripts.ScriptRepository) (*ScriptSearchUC, error) {
	if scriptS == nil {
		return nil, scripts.ErrInvalidScriptService
	}
	return &ScriptSearchUC{scriptS: scriptS}, nil
}

func (u *ScriptSearchUC) Search(ctx context.Context, namePart string) ([]ScriptDTO, error) {
	scripts, err := u.scriptS.SearchScripts(ctx, namePart)
	if err != nil {
		return nil, err
	}

	dto := make([]ScriptDTO, 0, len(scripts))
	for _, s := range scripts {
		dto = append(dto, ScriptToDTO(s))
	}
	return dto, nil
}
