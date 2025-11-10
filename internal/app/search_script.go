package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type SearchScriptsUC struct {
	scriptR scripts.ScriptRepository
	userP   scripts.UserProvider
	logger  *slog.Logger
}

func NewSearchScriptsUC(
	scriptR scripts.ScriptRepository,
	userP scripts.UserProvider,
	logger *slog.Logger,
) SearchScriptsUC {
	return SearchScriptsUC{scriptR: scriptR, userP: userP, logger: logger}
}

func (u *SearchScriptsUC) Search(ctx context.Context, userID uint32, substr string) ([]ScriptDTO, error) {
	u.logger.Info("searching scripts", "substr", substr)

	u.logger.Debug("getting user", "userID", userID, "ctx", ctx)
	user, err := u.userP.User(ctx, scripts.UserID(userID))
	u.logger.Debug("got user", "user", user, "err", err)
	if err != nil {
		u.logger.Error("failed to search script", "err", err.Error())
		return nil, err
	}

	u.logger.Debug("searching public scripts", "substr", substr)
	allScripts, err := u.scriptR.SearchPublicScripts(ctx, substr)
	u.logger.Debug("got scripts", "scripts count", len(allScripts), "err", err)
	if err != nil {
		u.logger.Error("failed to search script", "err", err.Error())
		return nil, err
	}

	u.logger.Debug("checking if user is admin", "is", user.IsAdmin())
	if !user.IsAdmin() {
		u.logger.Debug("searching user scripts", "substr", substr)
		userScripts, err := u.scriptR.SearchUserScripts(ctx, scripts.UserID(userID), substr)
		u.logger.Debug("got scripts", "scripts count", len(userScripts), "err", err)
		if err != nil {
			u.logger.Error("failed to search script", "err", err.Error())
			return nil, err
		}
		allScripts = append(allScripts, userScripts...)
	}

	u.logger.Debug("converting scripts to DTO", "scripts count", len(allScripts))
	dto := make([]ScriptDTO, 0, len(allScripts))
	for _, s := range allScripts {
		u.logger.Debug("converting script to DTO", "script", s)
		script, err := ScriptToDTO(s)
		u.logger.Debug("got DTO", "script", script, "err", err)
		if err != nil {
			u.logger.Error("failed to search script", "err", err.Error())
			return nil, err
		}
		dto = append(dto, script)
	}

	u.logger.Info("returning scripts", "scripts count", len(dto))
	return dto, nil
}
