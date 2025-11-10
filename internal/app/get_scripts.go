package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type GetScriptsUC struct {
	scriptR scripts.ScriptRepository
	userP   scripts.UserProvider
	logger  *slog.Logger
}

func NewGetScriptsUС(
	scriptR scripts.ScriptRepository,
	userP scripts.UserProvider,
	logger *slog.Logger,
) GetScriptsUC {
	return GetScriptsUC{scriptR: scriptR, userP: userP, logger: logger}
}

func (u *GetScriptsUC) Scripts(ctx context.Context, userID uint32) ([]ScriptDTO, error) {
	u.logger.Info("get scripts for user", "userID", userID)
	u.logger.Debug("get user", "userID", userID)
	user, err := u.userP.User(ctx, scripts.UserID(userID))
	u.logger.Debug("got user", "user", *user, "err", err)

	if err != nil {
		u.logger.Error("failed to get scripts for user", "err", err.Error())
		return nil, err
	}

	u.logger.Debug("get public scripts", "userID", userID)
	allScripts, err := u.scriptR.PublicScripts(ctx)
	u.logger.Debug("got scripts", "public scripts count", len(allScripts), "err", err)
	if err != nil {
		u.logger.Error("failed to get scripts for user", "err", err.Error())
		return nil, err
	}

	u.logger.Debug("is user admin", "is", user.IsAdmin())
	if !user.IsAdmin() {
		u.logger.Debug("user is not admin")
		u.logger.Debug("get user scripts", "userID", userID)
		userScripts, err := u.scriptR.UserScripts(ctx, scripts.UserID(userID))
		u.logger.Debug("got scripts", "user scripts count", len(userScripts), "err", err)
		if err != nil {
			u.logger.Error("failed to get scripts for user", "err", err.Error())
			return nil, err
		}
		allScripts = append(allScripts, userScripts...)
		u.logger.Debug("total scripts", "count", len(allScripts))
	}

	dto := make([]ScriptDTO, 0, len(allScripts))
	for _, s := range allScripts {
		u.logger.Debug("convert script to dto")
		script, err := ScriptToDTO(s)
		u.logger.Debug("got dto", "dto", script, "err", err)
		if err != nil {
			u.logger.Error("failed to get scripts for user", "err", err.Error())
			return nil, err
		}
		dto = append(dto, script)
	}

	u.logger.Info("return scripts", "count", len(dto))
	return dto, nil
}
