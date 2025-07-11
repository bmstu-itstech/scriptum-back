package scripts

import "time"

type Path = string
type ScriptID = uint32

type Visibility int

const (
	GLOBAL Visibility = iota
	PRIVATE
)

type Script struct {
	Fields     []Field
	path       Path
	owner      UserID
	visibility Visibility
	createdAt  time.Time
}

type PythonScript struct {
	Interpreter Path
}

func (s Script) Path() Path {
	return s.path
}

func (s Script) Owner() UserID {
	return s.owner
}

func (s Script) Visibility() Visibility {
	return s.visibility
}

func (s Script) CreatedAt() time.Time {
	return s.createdAt
}

func (s Script) Assemble(input Vector) Job {
	return Job{
		JobID:     0,
		UserID:    s.Owner(),
		In:        input,
		Command:   s.Path(),
		startedAt: time.Now(),
	}
}
