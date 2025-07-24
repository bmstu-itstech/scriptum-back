package scripts

import "fmt"

type ValueType struct {
	s string
}

var (
	IntegerType = ValueType{"integer"}
	RealType    = ValueType{"real"}
	ComplexType = ValueType{"complex"}
)

func (t ValueType) IsZero() bool {
	return t.s == ""
}

func (t ValueType) String() string {
	return t.s
}

func NewValueType(s string) (*ValueType, error) {
	switch s {
	case "integer":
		return &IntegerType, nil
	case "real":
		return &RealType, nil
	case "complex":
		return &ComplexType, nil
	default:
		return nil, fmt.Errorf(
			"%w: invalid ValueType: expected one of ['integer', 'real', 'complex'], got '%s']",
			ErrInvalidInput, s,
		)
	}
}
