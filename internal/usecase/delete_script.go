package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/bmstu-itstech/scriptum-back/internal/service"
)

type ScriptDeleteUC struct {
	service service.ScriptService
}

func NewScriptDeleteUC(service service.ScriptService) (*ScriptDeleteUC, error) {
	if service == nil {
		return nil, scripts.ErrInvalidScriptService
	}

	return &ScriptDeleteUC{service: service}, nil
}

func (u *ScriptDeleteUC) DeleteScript(ctx context.Context, scriptID scripts.ScriptID) error {
	err := u.service.DeleteScript(ctx, scriptID)
	// логи
	return err
}
