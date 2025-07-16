package scripts

import "context"

type ScriptRepository interface {
	CreateScript(ctx context.Context, script Script) (ScriptID, error)
	Script(ctx context.Context, scriptID ScriptID) (Script, error)
	DeleteScript(ctx context.Context, scriptID ScriptID) error
	PublicScripts(ctx context.Context) ([]Script, error)
	UserScripts(ctx context.Context, userID UserID) ([]Script, error)
	SearchPublicScripts(ctx context.Context, substr string) ([]Script, error)
	SearchUserScripts(ctx context.Context, userID UserID, substr string) ([]Script, error)
	UpdateScript(ctx context.Context, script Script) error
}
