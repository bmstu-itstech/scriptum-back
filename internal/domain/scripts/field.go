package scripts

import (
	"fmt"
)

const FieldNameMaxLen = 80
const FieldDescriptionMaxLen = 512
const FieldUnitMaxLen = 40

type Field struct {
	typ_ ValueType
	name string // 0 <  len(name) <= FieldNameMaxLen
	desc string // 0 <= len(desc) <= FieldDescriptionMaxLen
	unit string // 0 <  len(unit) <= FieldUnitMaxLen
}

func NewField(typ ValueType, name string, desc string, unit string) (*Field, error) {
	if typ.IsZero() {
		// Это может случиться тогда и только тогда, когда typ был создан не через конструктор, то есть
		// это ошибка программиста
		return nil, fmt.Errorf("typ is empty")
	}

	if name == "" {
		return nil, fmt.Errorf("%w: invalid Field: expected not empty name len", ErrInvalidInput)
	}

	if len(name) > FieldNameMaxLen {
		return nil, fmt.Errorf(
			"%w: invalid Field: expected len(name) <= %d, got len(name) = %d",
			ErrInvalidInput, FieldNameMaxLen, len(name),
		)
	}

	if len(desc) > FieldDescriptionMaxLen {
		return nil, fmt.Errorf(
			"%w: invalid Field: expected len(desc) <= %d, got len(desc) = %d",
			ErrInvalidInput, FieldDescriptionMaxLen, len(desc),
		)
	}

	if unit == "" {
		return nil, fmt.Errorf("%w: invalid Field: expected not empty unit", ErrInvalidInput)
	}

	if len(unit) > FieldUnitMaxLen {
		return nil, fmt.Errorf(
			"%w: invalid Field: expected len(unit) <= %d, got len(unit) = %d",
			ErrInvalidInput, FieldUnitMaxLen, len(unit),
		)
	}

	return &Field{
		typ_: typ,
		name: name,
		desc: desc,
		unit: unit,
	}, nil
}

func (f *Field) ValueType() ValueType {
	return f.typ_
}

func (f *Field) Name() string {
	return f.name
}

func (f *Field) Description() string {
	return f.desc
}

func (f *Field) Unit() string {
	return f.unit
}
