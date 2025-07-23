package scripts

import "fmt"

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

func NewType(s string) (*Type, error) {
	switch s {
	case "integer":
		return &IntegerType, nil
	case "real":
		return &RealType, nil
	case "complex":
		return &ComplexType, nil
	default:
		return nil, fmt.Errorf("type: expected one of [integer, real, complex], got %q: %w", s, ErrTypeInvalid)
	}
}
