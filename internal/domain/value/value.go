package value

import (
	"fmt"
	"strconv"

	"github.com/bmstu-itstech/scriptum-back/internal/domain"
)

type Value struct {
	t Type
	s string
}

func NewIntegerValue(s string) (Value, error) {
	err := validateInteger(s)
	if err != nil {
		return Value{}, err
	}
	return Value{
		t: IntegerValueType,
		s: s,
	}, nil
}

func MustNewIntegerValue(s string) Value {
	v, err := NewValue(IntegerValueType, s)
	if err != nil {
		panic(err)
	}
	return v
}

func NewRealValue(s string) (Value, error) {
	err := validateReal(s)
	if err != nil {
		return Value{}, err
	}
	return Value{
		t: RealValueType,
		s: s,
	}, nil
}

func MustNewRealValue(s string) Value {
	v, err := NewValue(RealValueType, s)
	if err != nil {
		panic(err)
	}
	return v
}

func NewStringValue(s string) Value {
	return Value{t: StringValueType, s: s}
}

func NewValue(t Type, s string) (Value, error) {
	switch t {
	case IntegerValueType:
		return NewIntegerValue(s)

	case RealValueType:
		return NewRealValue(s)

	case StringValueType:
		return NewStringValue(s), nil
	}
	return Value{}, domain.NewInvalidInputError(
		"value-type-invalid",
		fmt.Sprintf("invalid value type: %q", t.String()),
	)
}

func MustNewValue(t Type, s string) Value {
	v, err := NewValue(t, s)
	if err != nil {
		panic(err)
	}
	return v
}

func validateInteger(s string) error {
	_, err := strconv.ParseInt(s, 10, 64)
	return err
}

func validateReal(s string) error {
	_, err := strconv.ParseFloat(s, 64)
	return err
}

func (v Value) String() string {
	return v.s
}

func (v Value) Type() Type {
	return v.t
}
