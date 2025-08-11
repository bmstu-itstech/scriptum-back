package service

import (
	"context"
	"strings"
	"sync"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type MockScriptRepo struct {
	sync.RWMutex
	m      map[scripts.ScriptID]scripts.Script
	lastID scripts.ScriptID
}

func NewMockScriptRepository() (*MockScriptRepo, error) {
	return &MockScriptRepo{
		m: make(map[scripts.ScriptID]scripts.Script),
	}, nil
}

func (r *MockScriptRepo) Create(ctx context.Context, script *scripts.ScriptPrototype) (*scripts.Script, error) {
	r.Lock()
	defer r.Unlock()

	r.lastID++

	newScript, err := script.Build(r.lastID)
	if err != nil {
		return nil, err
	}

	r.m[r.lastID] = *newScript

	return newScript, nil
}

func (r *MockScriptRepo) Restore(ctx context.Context, script *scripts.Script) (*scripts.Script, error) {
	r.Lock()
	defer r.Unlock()

	r.lastID++

	newScript, err := script.Build(r.lastID)
	if err != nil {
		return nil, err
	}

	r.m[r.lastID] = *newScript

	return newScript, nil
}

func (r *MockScriptRepo) Update(ctx context.Context, script *scripts.Script) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[script.ID()]; !ok {
		return nil
	}
	r.m[script.ID()] = *script
	return nil

}

func (r *MockScriptRepo) Delete(ctx context.Context, id scripts.ScriptID) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[id]; !ok {
		return nil
	}
	delete(r.m, id)
	return nil
}

func (r *MockScriptRepo) Script(ctx context.Context, id scripts.ScriptID) (scripts.Script, error) {
	r.RLock()
	defer r.RUnlock()

	script, ok := r.m[id]
	if !ok {
		return scripts.Script{}, nil
	}
	return script, nil
}

func (r *MockScriptRepo) UserScripts(ctx context.Context, userID scripts.UserID) ([]scripts.Script, error) {
	r.RLock()
	defer r.RUnlock()
	scriptArray := make([]scripts.Script, 0, len(r.m))
	for _, s := range r.m {
		if s.OwnerID() == userID {
			scriptArray = append(scriptArray, s)
		}
	}
	if len(scriptArray) == 0 {
		return nil, nil
	}
	return scriptArray, nil
}

func (r *MockScriptRepo) PublicScripts(ctx context.Context) ([]scripts.Script, error) {
	r.RLock()
	defer r.RUnlock()
	scriptArray := make([]scripts.Script, 0, len(r.m))
	for _, s := range r.m {
		if s.IsPublic() {
			scriptArray = append(scriptArray, s)
		}
	}
	if len(scriptArray) == 0 {
		return nil, nil
	}
	return scriptArray, nil
}

func (r *MockScriptRepo) SearchPublicScripts(ctx context.Context, substr string) ([]scripts.Script, error) {
	r.RLock()
	defer r.RUnlock()
	scriptArray := make([]scripts.Script, 0, len(r.m))
	for _, s := range r.m {
		if s.IsPublic() && strings.Contains(s.Name(), substr) {
			scriptArray = append(scriptArray, s)
		}
	}
	if len(scriptArray) == 0 {
		return nil, nil
	}
	return scriptArray, nil
}

func (r *MockScriptRepo) SearchUserScripts(ctx context.Context, userID scripts.UserID, substr string) ([]scripts.Script, error) {

	r.RLock()
	defer r.RUnlock()
	scriptArray := make([]scripts.Script, 0, len(r.m))
	for _, s := range r.m {
		if s.OwnerID() == userID && strings.Contains(s.Name(), substr) {
			scriptArray = append(scriptArray, s)
		}
	}
	if len(scriptArray) == 0 {
		return nil, nil
	}
	return scriptArray, nil
}
