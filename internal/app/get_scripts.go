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
	user, err := u.userP.User(ctx, scripts.UserID(userID))
	if err != nil {
		u.logger.Error("failed to get scripts for user", "err", err.Error())
		return nil, err
	}

	allScripts, err := u.scriptR.PublicScripts(ctx)
	if err != nil {
		u.logger.Error("failed to get scripts for user", "err", err.Error())
		return nil, err
	}

	if !user.IsAdmin() {
		userScripts, err := u.scriptR.UserScripts(ctx, scripts.UserID(userID))
		if err != nil {
			u.logger.Error("failed to get scripts for user", "err", err.Error())
			return nil, err
		}
		allScripts = append(allScripts, userScripts...)
	}

	dto := make([]ScriptDTO, 0, len(allScripts))
	for _, s := range allScripts {
		script, err := ScriptToDTO(s)
		if err != nil {
			u.logger.Error("failed to get scripts for user", "err", err.Error())
			return nil, err
		}
		dto = append(dto, script)
	}

	return dto, nil
}
