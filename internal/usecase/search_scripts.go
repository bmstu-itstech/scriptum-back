package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptSearchUC struct {
	scriptS scripts.ScriptRepository
	userS   scripts.UserRepository
}

func NewScriptSearchUC(scriptS scripts.ScriptRepository, userS scripts.UserRepository) (*ScriptSearchUC, error) {
	if scriptS == nil {
		return nil, scripts.ErrInvalidScriptService
	}
	if userS == nil {
		return nil, scripts.ErrInvalidUserService
	}
	return &ScriptSearchUC{scriptS: scriptS, userS: userS}, nil
}

func (u *ScriptSearchUC) Search(ctx context.Context, userID uint32, substr string) ([]ScriptDTO, error) {
	allScripts, err := u.scriptS.SearchPublicScripts(ctx, substr)
	if err != nil {
		return nil, err
	}

	user, err := u.userS.User(ctx, scripts.UserID(userID))
	if err != nil {
		return nil, err
	}

	if !user.IsAdmin() {
		userScripts, err := u.scriptS.SearchUserScripts(ctx, scripts.UserID(userID), substr)
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
