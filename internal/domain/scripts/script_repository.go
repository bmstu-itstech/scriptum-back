package scripts

import "context"

type ScriptRepository interface {
	CreateScript(ctx context.Context, script Script) (ScriptID, error)
	Script(ctx context.Context, scriptID ScriptID) (Script, error)
	DeleteScript(ctx context.Context, scriptID ScriptID) error
	PublicScripts(ctx context.Context) ([]Script, error)
	UserScripts(ctx context.Context, userID UserID) ([]Script, error)
	SearchScripts(ctx context.Context, scriptNamePart string) ([]Script, error)
}
