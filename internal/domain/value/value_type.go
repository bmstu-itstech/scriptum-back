package value

type Type struct {
	s string
}

var (
	IntegerValueType = Type{"integer"}
	RealValueType    = Type{"real"}
	StringValueType  = Type{"string"}
)

func (t Type) IsZero() bool {
	return t == Type{}
}

func (t Type) String() string {
	return t.s
}
