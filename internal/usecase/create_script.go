package usecase

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
		panic(scripts.ErrInvalidScriptRepository)
	}
	if userR == nil {
		panic(scripts.ErrInvalidUserRepository)
	}
	if manager == nil {
		panic(scripts.ErrInvalidManagerService)
	}
	return &ScriptCreateUC{
		scriptR: scriptR,
		userR:   userR,
		manager: manager,
	}, nil
}

type ScriptCreateInput struct {
	UserID            uint32
	ScriptName        string
	ScriptDescription string
	File              FileDTO
	InFields          []FieldDTO
	OutFields         []FieldDTO
}

func (u *ScriptCreateUC) CreateScript(ctx context.Context, input ScriptCreateInput) (uint32, error) {
	user, err := u.userR.User(ctx, scripts.UserID(input.UserID))
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
		Owner:       input.UserID,
		Visibility:  string(vis),
	}

	script, err := DTOToScript(dto)
	if err != nil {
		// логируем ошибку
		return 0, err
	}

	scriptId, err := u.scriptR.Store(ctx, script)
	if err != nil {
		// логируем ошибку
		return 0, err
	}
	return uint32(scriptId), nil
}
