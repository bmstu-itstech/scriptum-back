package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type ScriptService interface {
	CreateScript(ctx context.Context, script scripts.Script) (scripts.ScriptID, error)
	GetScript(ctx context.Context, scriptID scripts.ScriptID) (scripts.Script, error)
	DeleteScript(ctx context.Context, scriptID scripts.ScriptID) error
	GetScripts(ctx context.Context) ([]scripts.Script, error)
	GetUserScripts(ctx context.Context, userID scripts.UserID) ([]scripts.Script, error)
}

type MockScriptRepo struct {
	sync.RWMutex
	m map[scripts.ScriptID]scripts.Script
}

func NewMockScriptRepository() (*MockScriptRepo, error) {
	return &MockScriptRepo{
		m: make(map[scripts.ScriptID]scripts.Script),
	}, nil
}

func (r *MockScriptRepo) GetScript(_ context.Context, scriptID scripts.ScriptID) (scripts.Script, error) {
	r.RLock()
	defer r.RUnlock()

	script, ok := r.m[scriptID]
	if !ok {
		return scripts.Script{}, fmt.Errorf("user not found: %d", scriptID)
	}
	return script, nil
}
