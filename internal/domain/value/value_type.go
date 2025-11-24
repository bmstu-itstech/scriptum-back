package value

import (
	"fmt"

	"github.com/bmstu-itstech/scriptum-back/internal/domain"
)

type Type struct {
	s string
}

var (
	IntegerValueType = Type{"integer"}
	RealValueType    = Type{"real"}
	StringValueType  = Type{"string"}
)

func TypeFromString(s string) (Type, error) {
	switch s {
	case "integer":
		return IntegerValueType, nil
	case "real":
		return RealValueType, nil
	case "string":
		return StringValueType, nil
	}
	return Type{}, domain.NewInvalidInputError(
		"type-invalid",
		fmt.Sprintf("invalid value type: expected one of ['integer', 'real', 'string'], got %s", s),
	)
}

func (t Type) IsZero() bool {
	return t == Type{}
}

func (t Type) String() string {
	return t.s
}
