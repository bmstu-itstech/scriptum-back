package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetScriptsUC struct {
	scriptR scripts.ScriptRepository
	userR   scripts.UserRepository
}

func NewGetScriptsUС(scriptR scripts.ScriptRepository, userR scripts.UserRepository) (*GetScriptsUC, error) {
	if scriptR == nil {
		panic(scripts.ErrInvalidScriptRepository)
	}

	if userR == nil {
		panic(scripts.ErrInvalidUserRepository)
	}

	return &GetScriptsUC{scriptR: scriptR, userR: userR}, nil
}

func (u *GetScriptsUC) Scripts(ctx context.Context, userID uint32) ([]ScriptDTO, error) {
	user, err := u.userR.User(ctx, scripts.UserID(userID))
	if err != nil {
		return nil, err
	}

	allScripts, err := u.scriptR.PublicScripts(ctx)
	if err != nil {
		return nil, err
	}

	if !user.IsAdmin() {
		userScripts, err := u.scriptR.UserScripts(ctx, scripts.UserID(userID))
		if err != nil {
			return nil, err
		}
		allScripts = append(allScripts, userScripts...)
	}

	dto := make([]ScriptDTO, 0, len(allScripts))
	for _, s := range allScripts {
		dto = append(dto, ScriptToDTO(s))
	}

	return dto, nil
}
