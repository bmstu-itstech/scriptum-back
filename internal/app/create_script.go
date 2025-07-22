package app

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptCreateUC struct {
	scriptR scripts.ScriptRepository
	userR   scripts.UserRepository
	manager scripts.Manager
}

func NewScriptCreateUC(
	scriptR scripts.ScriptRepository,
	userR scripts.UserRepository,
	manager scripts.Manager,
) (*ScriptCreateUC, error) {
	if scriptR == nil {
		return nil, scripts.ErrInvalidScriptRepository
	}
	if userR == nil {
		return nil, scripts.ErrInvalidUserRepository
	}
	if manager == nil {
		return nil, scripts.ErrInvalidManagerService
	}
	return &ScriptCreateUC{
		scriptR: scriptR,
		userR:   userR,
		manager: manager,
	}, nil
}

func (u *ScriptCreateUC) CreateScript(ctx context.Context, userID uint32, input ScriptCreateDTO) (uint32, error) {
	user, err := u.userR.User(ctx, scripts.UserID(userID))
	if err != nil {
		return 0, err
	}

	var vis scripts.Visibility
	if user.IsAdmin() {
		vis = scripts.VisibilityGlobal
	} else {
		vis = scripts.VisibilityPrivate
	}

	file, err := DTOToFile(input.File)
	if err != nil {
		// логируем ошибку
		return 0, err
	}

	path, err := u.manager.Upload(ctx, file)
	if err != nil {
		return 0, err
	}

	dto := ScriptDTO{
		InFields:    input.InFields,
		OutFields:   input.OutFields,
		Name:        input.ScriptName,
		Description: input.ScriptDescription,
		Path:        path,
		Owner:       userID,
		Visibility:  string(vis),
	}

	script, err := DTOToScript(dto)
	if err != nil {
		// логируем ошибку
		return 0, err
	}

	scriptId, err := u.scriptR.StoreScript(ctx, script)
	if err != nil {
		// логируем ошибку
		return 0, err
	}
	return uint32(scriptId), nil
}
