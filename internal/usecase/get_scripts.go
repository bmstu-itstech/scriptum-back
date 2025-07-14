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

func (u *GetScriptsUC) ScriptS() service.ScriptService {
	return u.scriptS
}

func (u *GetScriptsUC) UserS() userspb.UserServiceClient {
	return u.userS
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
	user, err := u.UserS().GetUser(ctx, &userspb.GetUserRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	if user.Visibility() == scripts.VisibilityGlobal {
		return u.ScriptS().GetScripts(ctx)
	} else if user.Visibility() == scripts.VisibilityPrivate {
		return u.ScriptS().GetUserScripts(ctx, userID)
	} 
	return nil, scripts.ErrInvalidVisibility
}
