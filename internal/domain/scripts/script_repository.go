package scripts

import "context"

type ScriptRepository interface {
	Store(ctx context.Context, script Script) (ScriptID, error)
	Script(ctx context.Context, scriptID ScriptID) (Script, error)
	Delete(ctx context.Context, scriptID ScriptID) error
	UserScripts(ctx context.Context, userID UserID) ([]Script, error)
	SearchPublicScripts(ctx context.Context, substr string) ([]Script, error)
	SearchUserScripts(ctx context.Context, userID UserID, substr string) ([]Script, error)
	PublicScripts(ctx context.Context) ([]Script, error)
	UpdateScript(ctx context.Context, script Script) error
}
