package scripts

import "time"

type Path = string
type ScriptID = uint32

type Visibility int

const (
	GLOBAL Visibility = iota
	PRIVATE
)

type PythonScript struct {
	interpreter Path
}

func (p *PythonScript) Interpreter() Path {
	return p.interpreter
}

func NewPythonScript(interpreter Path) (*PythonScript, error) {
	if interpreter == "" {
		return nil, ErrInvalidInterpreter
	}
	return &PythonScript{
		interpreter: interpreter,
	}, nil
}

type Script struct {
	fields     []Field
	path       Path
	owner      UserID
	visibility Visibility
	createdAt  time.Time
}

func (s *Script) Fields() []Field {
	return s.fields
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

func NewScript(fields []Field, path Path, owner UserID, visibility Visibility) (*Script, error) {
	if len(fields) == 0 {
		return nil, ErrFieldsEmpty
	}
	if path == "" {
		return nil, ErrPathEmpty
	}

	return &Script{
		fields:     fields,
		path:       path,
		owner:      owner,
		visibility: visibility,
		createdAt:  time.Now(),
	}, nil
}

func (s *Script) Assemble(input Vector) (*Job, error) {
	return &Job{
		jobID:     0,
		userID:    s.Owner(),
		in:        input,
		command:   s.Path(),
		startedAt: time.Now(),
	}, nil
}
