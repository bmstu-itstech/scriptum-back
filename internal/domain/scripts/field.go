package scripts

import (
	"fmt"
	"strings"
)

type Field struct {
	fieldType Type
	name      string
	desc      string
	unit      string
}

const FieldNameMaxLength = 100
const FieldDescMaxLength = 500
const FieldUnitMaxLength = 20

func (f *Field) FieldType() Type {
	return f.fieldType
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

func NewField(fieldType Type, name, desc, unit string) (*Field, error) {
	if name == "" {
		return nil, fmt.Errorf("name: expected not empty string, got empty string  %w", ErrFieldInvalid)
	}
	if len(name) > FieldNameMaxLength {
		return nil, fmt.Errorf("name: expected string with len below %d, got string length %d %w", FieldNameMaxLength, len(name), ErrFieldInvalid)
	}
	if desc == "" {
		return nil, fmt.Errorf("desc: expected not empty string, got empty string  %w", ErrFieldInvalid)
	}
	if len(desc) > FieldDescMaxLength {
		return nil, fmt.Errorf("desc: expected string with len below %d, got string length %d %w", FieldDescMaxLength, len(desc), ErrFieldInvalid)
	}
	if unit == "" {
		return nil, fmt.Errorf("unit: expected not empty string, got empty string  %w", ErrFieldInvalid)
	}
	if len(desc) > FieldUnitMaxLength {
		return nil, fmt.Errorf("unit: expected string with len below %d, got string length %d %w", FieldUnitMaxLength, len(unit), ErrFieldInvalid)
	}

	return &Field{
		fieldType: fieldType,
		name:      name,
		desc:      desc,
		unit:      unit,
	}, nil
}

func ParseOutputValues(output string, fields []Field) ([]Value, error) {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) != len(fields) {
		return nil, ErrFieldCount
	}

	values := make([]Value, 0, len(fields))
	for i, line := range lines {
		val, err := NewValue(fields[i].FieldType().String(), line)
		if err != nil {
			return nil, err
		}
		values = append(values, val)
	}
	return values, nil
}
