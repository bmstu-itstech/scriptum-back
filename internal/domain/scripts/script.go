package scripts

import (
	"fmt"
	"time"
)

type Path = string
type ScriptID = uint32

type Visibility string

const ScriptNameMaxLength = 100
const ScriptPathMaxLength = 200
const ScriptDescriptionMaxLength = 500

const (
	VisibilityGlobal  Visibility = "global"
	VisibilityPrivate Visibility = "private"
)

type Script struct {
	id          ScriptID
	name        string
	description string
	inFields    []Field
	outFields   []Field
	path        Path
	owner       UserID
	visibility  Visibility
	createdAt   time.Time
}

func (s *Script) ID() ScriptID {
	return s.id
}

func (s *Script) InFields() []Field {
	return s.inFields
}

func (s *Script) OutFields() []Field {
	return s.outFields
}

func (s *Script) Path() Path {
	return s.path
}

func (s *Script) Owner() UserID {
	return s.owner
}

func (s *Script) Visibility() Visibility {
	return s.visibility
}

func (s *Script) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Script) Name() string {
	return s.name
}

func (s *Script) Description() string {
	return s.description
}

func IsGlobal(v Visibility) bool {
	switch v {
	case VisibilityGlobal:
		return true
	default:
		return false
	}
}

func NewScript(
	scriptID ScriptID,
	inFields []Field,
	outFields []Field,
	path Path,
	owner UserID,
	visibility Visibility,
	name string,
	description string,
) (*Script, error) {
	if len(inFields) == 0 {
		return nil, fmt.Errorf("inFields: expected non-empty slice, got length %d: %w", len(inFields), ErrScriptInvalid)
	}
	if len(outFields) == 0 {
		return nil, fmt.Errorf("outFields: expected non-empty slice, got length %d: %w", len(outFields), ErrScriptInvalid)
	}

	if path == "" {
		return nil, fmt.Errorf("path: expected non-empty string, got empty string: %w", ErrScriptInvalid)
	}
	if len(path) > ScriptPathMaxLength {
		return nil, fmt.Errorf("path: expected string with length ≤ %d, got length %d: %w", ScriptPathMaxLength, len(path), ErrScriptInvalid)
	}

	if name == "" {
		return nil, fmt.Errorf("name: expected non-empty string, got empty string: %w", ErrScriptInvalid)
	}
	if len(name) > ScriptNameMaxLength {
		return nil, fmt.Errorf("name: expected string with length ≤ %d, got length %d: %w", ScriptNameMaxLength, len(name), ErrScriptInvalid)
	}

	if description == "" {
		return nil, fmt.Errorf("description: expected non-empty string, got empty string: %w", ErrScriptInvalid)
	}
	if len(description) > ScriptDescriptionMaxLength {
		return nil, fmt.Errorf("description: expected string with length ≤ %d, got length %d: %w", ScriptDescriptionMaxLength, len(description), ErrScriptInvalid)
	}

	return &Script{
		id:          scriptID,
		inFields:    inFields,
		outFields:   outFields,
		name:        name,
		description: description,
		path:        path,
		owner:       owner,
		visibility:  visibility,
		createdAt:   time.Now(),
	}, nil
}

func NewScriptRead(scriptID ScriptID, inFields []Field, outFields []Field, path Path, owner UserID, visibility Visibility, name string, description string, createdAt time.Time) (*Script, error) {
	return &Script{
		id:          scriptID,
		inFields:    inFields,
		outFields:   outFields,
		name:        name,
		description: description,
		path:        path,
		owner:       owner,
		visibility:  visibility,
		createdAt:   createdAt,
	}, nil
}

func (s *Script) Assemble(input Vector, email Email, needToNotify bool) (*Job, error) {
	return NewJob(0, s.Owner(), input, s.Path(), time.Now(), s.InFields(), email, needToNotify)
}
