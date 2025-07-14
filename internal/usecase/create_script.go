package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/bmstu-itstech/scriptum-back/internal/service"
)

type ScriptCreateUC struct {
	service service.ScriptService
}

func (s *ScriptCreateUC) Service() service.ScriptService {
	return s.service
}

func NewScriptCreateUC(service service.ScriptService) (*ScriptCreateUC, error) {
	if service == nil {
		return nil, scripts.ErrInvalidScriptService
	}
	return &ScriptCreateUC{service: service}, nil
}

func (u *ScriptCreateUC) CreateScript(ctx context.Context, userID scripts.UserID, script scripts.Script) (scripts.ScriptID, error) {

	scriptId, err := u.Service().CreateScript(ctx, script)
	if err != nil {
		// логируем ошибку
		return 0, err
	}
	return scriptId, nil
}
