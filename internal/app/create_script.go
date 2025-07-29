package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptCreateUC struct {
	scriptR scripts.ScriptRepository
	userP   scripts.UserProvider
	manager scripts.FileManager
	logger  *slog.Logger
}

func NewScriptCreateUC(
	scriptR scripts.ScriptRepository,
	userP scripts.UserProvider,
	manager scripts.FileManager,
	logger *slog.Logger,
) ScriptCreateUC {
	return ScriptCreateUC{
		scriptR: scriptR,
		userP:   userP,
		manager: manager,
		logger:  logger,
	}
}

func (u *ScriptCreateUC) CreateScript(ctx context.Context, req ScriptCreateDTO) (int32, error) {
	user, err := u.userP.User(ctx, scripts.UserID(req.OwnerID))
	if err != nil {
		return 0, err
	}

	var vis scripts.Visibility
	if user.IsAdmin() {
		vis = scripts.VisibilityPublic
	} else {
		vis = scripts.VisibilityPrivate
	}

	input, err := DTOToFields(req.InFields)
	if err != nil {
		return 0, err
	}

	output, err := DTOToFields(req.OutFields)
	if err != nil {
		return 0, err
	}

	file, err := DTOToFile(req.File)
	if err != nil {
		return 0, err
	}

	url, err := u.manager.Save(ctx, file)
	if err != nil {
		return 0, err
	}

	proto, err := scripts.NewScriptPrototype(
		scripts.UserID(req.OwnerID), req.ScriptName, req.ScriptDescription, vis, input, output, url,
	)
	if err != nil {
		return 0, err
	}

	script, err := u.scriptR.Create(ctx, proto)
	if err != nil {
		return 0, err
	}

	return int32(script.ID()), nil
}
