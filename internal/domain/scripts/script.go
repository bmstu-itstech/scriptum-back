package scripts

import (
	"time"
)

type Path = string
type ScriptID = uint32

type Visibility string

const (
	VisibilityGlobal  Visibility = "global"
	VisibilityPrivate Visibility = "private"
)

// type PythonScript struct {
// 	interpreter Path
// }

// func (p *PythonScript) Interpreter() Path {
// 	return p.interpreter
// }

// func NewPythonScript(interpreter Path) (*PythonScript, error) {
// 	if interpreter == "" {
// 		return nil, ErrInvalidInterpreter
// 	}
// 	return &PythonScript{
// 		interpreter: interpreter,
// 	}, nil
// }

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

func NewScript(scriptID ScriptID, inFields []Field, outFields []Field, path Path, owner UserID, visibility Visibility, name string, description string) (*Script, error) {
	if (len(inFields) == 0) || (len(outFields) == 0) {
		return nil, ErrFieldsEmpty
	}

	if path == "" {
		return nil, ErrPathEmpty
	}

	if name == "" {
		return nil, ErrNameEmpty
	}

	if description == "" {
		return nil, ErrDescriptionEmpty
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

func (s *Script) Assemble(input Vector) (*Job, error) {
	return NewJob(0, s.Owner(), input, s.Path(), time.Now())
}
