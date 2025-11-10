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
	u.logger.Debug("create script debug", "req", req, "ctx", ctx)
	user, err := u.userP.User(ctx, scripts.UserID(req.OwnerID))
	u.logger.Debug("user provider user", "user", user, "err", err)
	if err != nil {
		u.logger.Error("failed to get user", "err", err.Error())
		return 0, err
	}

	var vis scripts.Visibility
	if user.IsAdmin() {
		vis = scripts.VisibilityPublic
	} else {
		vis = scripts.VisibilityPrivate
	}
	u.logger.Debug("user is admin", "isAdmin", vis)

	input, err := DTOToFields(req.InFields)
	u.logger.Debug("input fields", "input", input, "err", err)
	if err != nil {
		u.logger.Error("failed to convert input fields", "err", err.Error())
		return 0, err
	}

	output, err := DTOToFields(req.OutFields)
	u.logger.Debug("output fields", "output", output, "err", err)
	if err != nil {
		u.logger.Error("failed to convert output fields", "err", err.Error())
		return 0, err
	}

	extraFileIDs := make([]scripts.FileID, len(req.ExtraFileIDs))
	for i, id := range req.ExtraFileIDs {
		extraFileIDs[i] = scripts.FileID(id)
	}
	u.logger.Debug("extra file ids", "extraFileIDs", extraFileIDs)

	u.logger.Debug("script prototype creating")
	proto, err := scripts.NewScriptPrototype(
		scripts.UserID(req.OwnerID),
		req.ScriptName,
		req.ScriptDescription,
		vis,
		scripts.PythonVersion(req.PythonVersion),
		input,
		output,
		scripts.FileID(req.MainFileID),
		extraFileIDs,
	)
	u.logger.Debug("script prototype", "proto", proto, "err", err)

	if err != nil {
		u.logger.Error("failed to create script prototype", "err", err.Error())
		return 0, err
	}

	u.logger.Debug("script repository creating")
	script, err := u.scriptR.Create(ctx, proto)
	u.logger.Debug("script repository", "script", script, "err", err)
	if err != nil {
		u.logger.Error("failed to create script", "err", err.Error())
		return 0, err
	}

	u.logger.Info("script created", "script", script)
	return int32(script.ID()), nil
}
