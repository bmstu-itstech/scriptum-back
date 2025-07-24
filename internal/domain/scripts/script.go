package scripts

import (
	"fmt"
	"time"
)

const ScriptNameMaxLen = 64
const ScriptDescriptionMaxLen = 256

type ScriptID int32

type Visibility struct {
	s string
}

var VisibilityPublic = Visibility{"public"}
var VisibilityPrivate = Visibility{"private"}

func (v Visibility) String() string {
	return v.s
}

func (v Visibility) IsZero() bool {
	return v.s == ""
}

type ScriptPrototype struct {
	ownerID UserID     // ownerID != 0
	name    string     // 0 <  len(name) < ScriptNameMaxLen
	desc    string     // 0 <= len(desc) < ScriptDescriptionMaxLen
	vis     Visibility // !vis.IsZero()
	input   []Field    // len(input) > 0
	output  []Field    // len(output) > 0
	url     URL        // len(url) > 0
}

func NewScriptPrototype(
	ownerID UserID,
	name string,
	desc string,
	visibility Visibility,
	input []Field,
	output []Field,
	url URL,
) (*ScriptPrototype, error) {
	if ownerID == 0 {
		// Ошибка программиста
		return nil, fmt.Errorf("empty ownerID")
	}

	if name == "" {
		return nil, fmt.Errorf("%w: invalid Script: expected not empty name", ErrInvalidInput)
	}

	if len(name) > ScriptNameMaxLen {
		return nil, fmt.Errorf(
			"%w: invalid Script: expected len(name) < %d, got len(name) = %d",
			ErrInvalidInput, FieldNameMaxLen, len(name),
		)
	}

	if len(desc) > ScriptDescriptionMaxLen {
		return nil, fmt.Errorf(
			"%w: invalid Script: expected len(desc) < %d, got len(desc) = %d",
			ErrInvalidInput, FieldDescriptionMaxLen, len(desc),
		)
	}

	if visibility.IsZero() {
		// Visibility не является пользовательским вводом.
		// Пустой visibility есть ошибка программиста.
		return nil, fmt.Errorf("empty visibility")
	}

	if len(input) == 0 {
		return nil, fmt.Errorf("%w: invalid Script: expected at least one input field", ErrInvalidInput)
	}

	if len(output) == 0 {
		return nil, fmt.Errorf("%w: invalid Script: expected at least one output field", ErrInvalidInput)
	}

	if len(url) == 0 {
		return nil, fmt.Errorf("%w: invalid Script: expected not empty URL", ErrInvalidInput)
	}

	return &ScriptPrototype{
		ownerID: ownerID,
		name:    name,
		desc:    desc,
		vis:     visibility,
		input:   input[:],
		output:  output[:],
		url:     url,
	}, nil
}

func (s *ScriptPrototype) OwnerID() UserID {
	return s.ownerID
}

func (s *ScriptPrototype) Name() string {
	return s.name
}

func (s *ScriptPrototype) Desc() string {
	return s.desc
}

func (s *ScriptPrototype) Visibility() Visibility {
	return s.vis
}

func (s *ScriptPrototype) Input() []Field {
	return s.input[:]
}

func (s *ScriptPrototype) Output() []Field {
	return s.output[:]
}

func (s *ScriptPrototype) URL() URL {
	return s.url
}

func (s *ScriptPrototype) IsAvailableFor(userID UserID) bool {
	if s.vis == VisibilityPublic {
		return true
	}
	return s.ownerID == userID
}

func (s *ScriptPrototype) Build(id ScriptID) (*Script, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid Script: expected positive id, got %d", ErrInvalidInput, id)
	}

	return &Script{
		ScriptPrototype: *s,
		id:              id,
		createdAt:       time.Now(),
	}, nil
}

type Script struct {
	ScriptPrototype
	id        ScriptID
	createdAt time.Time
}

func (s *Script) ID() ScriptID {
	return s.id
}

func (s *Script) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Script) Assemble(by UserID, input []Value) (*JobPrototype, error) {
	if len(s.ScriptPrototype.input) != len(input) {
		return nil, fmt.Errorf(
			"%w: failed to assemble job: expected %d input values, got %d",
			ErrInvalidInput, len(s.ScriptPrototype.input), len(input),
		)
	}

	for i, field := range s.ScriptPrototype.input {
		value := input[i]
		if field.ValueType() != value.Type() {
			return nil, fmt.Errorf(
				"%w: failed to assemble job: expected type of input[i] is %s, got %s",
				ErrInvalidInput, field.ValueType(), value.Type(),
			)
		}
	}

	return &JobPrototype{
		ownerID:   by,
		input:     input,
		createdAt: time.Now(),
	}, nil
}
