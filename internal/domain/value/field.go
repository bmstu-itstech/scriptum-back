package value

import (
	"errors"
	"fmt"

	"github.com/bmstu-itstech/scriptum-back/internal/domain"
)

type Field struct {
	t    Type
	name string
	desc *string
	unit *string
}

func NewField(t Type, name string, desc *string, unit *string) (Field, error) {
	if t.IsZero() {
		return Field{}, errors.New("field type is zero")
	}

	if name == "" {
		return Field{}, domain.NewInvalidInputError("field-empty-name", "expected not empty field name")
	}

	if desc != nil && *desc == "" {
		return Field{}, errors.New("field description is not nil but empty")
	}

	if unit != nil && *unit == "" {
		return Field{}, errors.New("field unit is not nil but empty")
	}

	return Field{
		t:    t,
		name: name,
		desc: desc,
		unit: unit,
	}, nil
}

func (f Field) Validate(v Value) error {
	if f.t != v.t {
		return domain.NewInvalidInputError(
			"field-mismatch",
			fmt.Sprintf("field type mismatch: expected '%s', got '%s'", f.t.String(), v.t.String()),
		)
	}
	return nil
}

func (f Field) Type() Type {
	return f.t
}

func (f Field) Name() string {
	return f.name
}

func (f Field) Desc() *string {
	return f.desc
}

func (f Field) Unit() *string {
	return f.unit
}
