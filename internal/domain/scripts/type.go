package scripts

type Type struct {
	s string
}

var (
	IntegerType = Type{"integer"}
	RealType    = Type{"real"}
	ComplexType = Type{"complex"}
)

func (t Type) String() string {
	return t.s
}

func (t Type) IsInteger() bool {
	return t == IntegerType
}

func (t Type) IsReal() bool {
	return t == RealType
}

func (t Type) IsComplex() bool {
	return t == ComplexType
}

func NewType(s string) (Type, error) {
	switch s {
	case "integer":
		return IntegerType, nil
	case "real":
		return RealType, nil
	case "complex":
		return ComplexType, nil
	default:
		return Type{}, ErrInvalidType
	}
}
