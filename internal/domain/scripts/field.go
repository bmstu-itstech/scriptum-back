package scripts

import (
	"strings"
)

type Field struct {
	fieldType Type
	name      string
	desc      string
	unit      string
}

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
		return nil, ErrFieldNameEmpty
	}
	if desc == "" {
		return nil, ErrFieldDescEmpty
	}
	if unit == "" {
		return nil, ErrFieldUnitEmpty
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
		return nil, ErrFieldCnt
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
