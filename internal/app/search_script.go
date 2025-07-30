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
) GetScriptsUC {
	return GetScriptsUC{scriptR: scriptR, userP: userP, logger: logger}
}

func (u *SearchScriptsUC) Search(ctx context.Context, userID uint32, substr string) ([]ScriptDTO, error) {
	user, err := u.userP.User(ctx, scripts.UserID(userID))
	if err != nil {
		return nil, err
	}

	allScripts, err := u.scriptR.SearchPublicScripts(ctx, substr)
	if err != nil {
		return nil, err
	}

	if !user.IsAdmin() {
		userScripts, err := u.scriptR.SearchUserScripts(ctx, scripts.UserID(userID), substr)
		if err != nil {
			return nil, err
		}
		allScripts = append(allScripts, userScripts...)
	}

	dto := make([]ScriptDTO, 0, len(allScripts))
	for _, s := range allScripts {
		script, err := ScriptToDTO(s)
		if err != nil {
			return nil, err
		}
		dto = append(dto, script)
	}

	return dto, nil
}
