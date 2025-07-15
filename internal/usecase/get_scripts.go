package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/bmstu-itstech/scriptum-back/internal/service"
)

type GetScriptsUC struct {
	scriptS scripts.ScriptRepository
	userS   service.UserServiceClient
}

func NewGetScriptsUС(scriptS scripts.ScriptRepository, userS userspb.UserServiceClient) (*GetScriptsUC, error) {
	if scriptS == nil {
		return nil, scripts.ErrInvalidScriptService
	}

	if userS == nil {
		return nil, scripts.ErrInvalidUserService
	}

	return &GetScriptsUC{scriptS: scriptS, userS: userS}, nil
}

func (u *GetScriptsUC) Scripts(ctx context.Context, userID uint32) ([]ScriptDTO, error) {
	var err error
	var gotScripts []scripts.Script
	var user scripts.User

	user, err = u.userS.User(ctx, &userspb.GetUserRequest{UserId: userID})
	if err != nil {
		return nil, err
	}

	gotScripts, err = u.scriptS.PublicScripts(ctx)
	if err != nil {
		return nil, err
	}

	if !user.IsAdmin() {
		userScripts, err := u.scriptS.UserScripts(ctx, scripts.UserID(userID))
		if err != nil {
			return nil, err
		}
		gotScripts = append(gotScripts, userScripts...)
	}

	scriptsOut := make([]ScriptDTO, 0, len(gotScripts))
	for _, script := range gotScripts {
		scriptsOut = append(scriptsOut, ScriptToDTO(script))
	}

	return scriptsOut, nil
}
