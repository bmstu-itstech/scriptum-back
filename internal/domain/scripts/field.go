package scripts

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
