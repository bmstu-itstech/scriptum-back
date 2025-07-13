package service

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type MockScriptRepo struct {
	sync.RWMutex
	m      map[scripts.ScriptID]scripts.Script
	nextID uint32
}

func NewMockScriptRepository() (*MockScriptRepo, error) {
	return &MockScriptRepo{
		m: make(map[scripts.ScriptID]scripts.Script),
	}, nil
}

func (r *MockScriptRepo) GetScript(scriptID scripts.ScriptID) (scripts.Script, error) {
	r.RLock()
	defer r.RUnlock()

	script, ok := r.m[scriptID]
	if !ok {
		return scripts.Script{}, fmt.Errorf("script not found: %d", scriptID)
	}
	return script, nil
}

func (r *MockScriptRepo) StoreScript(script scripts.Script) (scripts.ScriptID, error) {
	r.Lock()
	defer r.Unlock()

	newID := scripts.ScriptID(atomic.AddUint32(&r.nextID, 1))

	r.m[newID] = script
	return newID, nil
}
