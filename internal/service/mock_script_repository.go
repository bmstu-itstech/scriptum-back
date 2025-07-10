package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type MockScriptRepository interface {
	GetScript(scriptID scripts.ScriptID) (scripts.Script, error)
}

type MockScriptRepo struct {
	context context.Context
	sync.RWMutex
	m map[scripts.ScriptID]scripts.Script
}

func NewMockScriptRepository() *MockScriptRepo {
	return &MockScriptRepo{
		context: context.Background(),
		m:       make(map[scripts.ScriptID]scripts.Script),
	}
}

func (r *MockScriptRepo) GetScript(scriptID scripts.ScriptID) (scripts.Script, error) {
	r.RLock()
	defer r.RUnlock()

	script, ok := r.m[scriptID]
	if !ok {
		return scripts.Script{}, fmt.Errorf("user not found: %d", scriptID)
	}
	return script, nil
}
