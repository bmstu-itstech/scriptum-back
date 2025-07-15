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

func (u *GetScriptsUC) Scripts(ctx context.Context, userID uint32) ([]UseCaseScript, error) {
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
		userScripts, err := u.scriptS.UserScripts(ctx, userID)
		if err != nil {
			return nil, err
		}
		gotScripts = append(gotScripts, userScripts...)
	}

	scriptsOut := make([]UseCaseScript, len(gotScripts))
	for _, script := range gotScripts {
		ucFields := make([]UseCaseField, len(script.Fields()))
		for _, field := range script.Fields() {
			sType := field.FieldType()
			if err != nil {
				return nil, err
			}
			ucFields = append(ucFields, UseCaseField{
				Type:        sType.String(),
				Name:        field.Name(),
				Description: field.Description(),
				Unit:        field.Unit(),
			})
		}
		scriptsOut = append(scriptsOut, UseCaseScript{
			Fields:     ucFields,
			Path:       script.Path(),
			Owner:      int64(script.Owner()),
			Visibility: string(script.Visibility()),
			CreatedAt:  script.CreatedAt(),
		})
	}
	return scriptsOut, nil
}
