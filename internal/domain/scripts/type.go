package scripts

type TypeValue int

const (
	INTEGER TypeValue = iota
	REAL
	COMPLEX
)

type Type struct {
	value TypeValue
}

func (t *Type) Value() TypeValue {
	return t.value
}

func NewType(value TypeValue) (*Type, error) {
	return &Type{value: value}, nil
}
