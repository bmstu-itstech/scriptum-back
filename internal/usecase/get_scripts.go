package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetScriptsUC struct {
	scriptS scripts.ScriptRepository
	userS   scripts.UserRepository
}

func NewGetScriptsUС(scriptS scripts.ScriptRepository, userS scripts.UserRepository) (*GetScriptsUC, error) {
	if scriptS == nil {
		return nil, scripts.ErrInvalidScriptService
	}

	if userS == nil {
		return nil, scripts.ErrInvalidUserService
	}

	return &GetScriptsUC{scriptS: scriptS, userS: userS}, nil
}

func (u *GetScriptsUC) Scripts(ctx context.Context, userID uint32) ([]ScriptDTO, error) {
	user, err := u.userS.User(ctx, scripts.UserID(userID))
	if err != nil {
		return nil, err
	}

	allScripts, err := u.scriptS.GetPublicScripts(ctx)
	if err != nil {
		return nil, err
	}

	if !user.IsAdmin() {
		userScripts, err := u.scriptS.GetUserScripts(ctx, scripts.UserID(userID))
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
