package app

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptSearchUC struct {
	scriptR scripts.ScriptRepository
	userR   scripts.UserRepository
}

func NewScriptSearchUC(scriptR scripts.ScriptRepository, userR scripts.UserRepository) (*ScriptSearchUC, error) {
	if scriptR == nil {
		return nil, scripts.ErrInvalidScriptRepository
	}
	if userR == nil {
		return nil, scripts.ErrInvalidUserRepository
	}
	return &ScriptSearchUC{scriptR: scriptR, userR: userR}, nil
}

func (u *ScriptSearchUC) Search(ctx context.Context, userID uint32, substr string) ([]ScriptDTO, error) {
	user, err := u.userR.User(ctx, scripts.UserID(userID))
	if err != nil {
		return nil, err
	}

	allScripts, err := u.scriptR.SearchPublicScripts(ctx, substr)
	if err != nil {
		return nil, err
	}

	if !user.IsAdmin() {
		userScripts, err := u.scriptR.SearchUserScripts(ctx, scripts.UserID(userID), substr)
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
