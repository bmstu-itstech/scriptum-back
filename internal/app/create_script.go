package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptCreateUC struct {
	scriptR scripts.ScriptRepository
	fileR   scripts.FileRepository
	userP   scripts.UserProvider
	manager scripts.FileManager
	logger  *slog.Logger
}

func NewScriptCreateUC(
	scriptR scripts.ScriptRepository,
	userP scripts.UserProvider,
	fileR scripts.FileRepository,
	manager scripts.FileManager,
	logger *slog.Logger,
) ScriptCreateUC {
	return ScriptCreateUC{
		scriptR: scriptR,
		userP:   userP,
		fileR:   fileR,
		manager: manager,
		logger:  logger,
	}
}

func (u *ScriptCreateUC) CreateScript(ctx context.Context, req ScriptCreateDTO) (int32, error) {
	u.logger.Info("create script", "req", req)
	user, err := u.userP.User(ctx, scripts.UserID(req.OwnerID))
	if err != nil {
		u.logger.Error("failed to get user", "err", err)
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
		u.logger.Error("failed to convert input fields", "err", err)
		return 0, err
	}

	output, err := DTOToFields(req.OutFields)
	if err != nil {
		u.logger.Error("failed to convert output fields", "err", err)
		return 0, err
	}

	file, err := u.fileR.File(ctx, scripts.FileID(req.FileID))
	if err != nil {
		u.logger.Error("failed to save file", "err", err)
		return 0, err
	}

	proto, err := scripts.NewScriptPrototype(
		scripts.UserID(req.OwnerID), req.ScriptName, req.ScriptDescription, vis, input, output, *file,
	)
	if err != nil {
		u.logger.Error("failed to create script prototype", "err", err)
		return 0, err
	}

	script, err := u.scriptR.Create(ctx, proto)
	if err != nil {
		u.logger.Error("failed to create script", "err", err)
		return 0, err
	}

	u.logger.Info("script created", "script", script)
	return int32(script.ID()), nil
}
