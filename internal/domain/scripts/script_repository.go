package scripts

import (
	"context"
	"errors"
)

var ErrScriptNotFound = errors.New("script not found")

type ScriptRepository interface {
	Create(ctx context.Context, script *ScriptPrototype) (*Script, error)
	Restore(ctx context.Context, script *Script) (*Script, error)
	Update(ctx context.Context, script *Script) error
	Delete(ctx context.Context, id ScriptID) error

	Script(ctx context.Context, id ScriptID) (Script, error)
	UserScripts(ctx context.Context, userID UserID) ([]Script, error)
	PublicScripts(ctx context.Context) ([]Script, error)

	SearchPublicScripts(ctx context.Context, substr string) ([]Script, error)
	SearchUserScripts(ctx context.Context, userID UserID, substr string) ([]Script, error)
}
