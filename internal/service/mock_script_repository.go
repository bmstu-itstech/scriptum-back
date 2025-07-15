package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type MockScriptRepo struct {
	sync.RWMutex
	m map[scripts.ScriptID]scripts.Script
}

func NewMockScriptRepository() (*MockScriptRepo, error) {
	return &MockScriptRepo{
		m: make(map[scripts.ScriptID]scripts.Script),
	}, nil
}

func (r *MockScriptRepo) Script(_ context.Context, scriptID scripts.ScriptID) (scripts.Script, error) {
	r.RLock()
	defer r.RUnlock()

	script, ok := r.m[scriptID]
	if !ok {
		return scripts.Script{}, fmt.Errorf("user not found: %d", scriptID)
	}
	return script, nil
}
