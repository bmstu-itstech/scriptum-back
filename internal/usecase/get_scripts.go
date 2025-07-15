package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/bmstu-itstech/scriptum-back/internal/service"
)

type GetScriptsUC struct {
	scriptS service.ScriptService
	userS   userspb.UserServiceClient
}

func NewGetScriptsUС(scriptS service.ScriptService, userS userspb.UserServiceClient) (*GetScriptsUC, error) {
	if scriptS == nil {
		return nil, scripts.ErrInvalidScriptService
	}

	if userS == nil {
		return nil, scripts.ErrInvalidUserService
	}

	return &GetScriptsUC{scriptS: scriptS, userS: userS}, nil
}

func (u *GetScriptsUC) GetScripts(ctx context.Context, userID scripts.UserID) ([]scripts.Script, error) {
	user, err := u.userS.GetUser(ctx, &userspb.GetUserRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	if user.Visibility() == scripts.VisibilityGlobal {
		return u.scriptS.GetScripts(ctx)
	} else if user.Visibility() == scripts.VisibilityPrivate {
		return u.scriptS.GetUserScripts(ctx, userID)
	} 
	return nil, scripts.ErrInvalidVisibility
}
