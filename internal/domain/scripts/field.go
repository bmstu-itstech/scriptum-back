package scripts

type Field struct {
	FieldType Type
	name      string
	desc      string
	unit      string
}

func (f Field) Name() string {
	return f.name
}

func (f Field) Description() string {
	return f.desc
}

func (f Field) Unit() string {
	return f.unit
}
