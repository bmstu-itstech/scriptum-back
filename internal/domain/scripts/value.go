package scripts

import (
	"fmt"
	"strconv"
)

type Value interface {
	VariableType() Type
	String() string
}

type Complex struct {
	data complex64
}

func NewComplex(data complex64) (*Complex, error) {
	return &Complex{data: data}, nil
}

func (c *Complex) Data() complex64 {
	return c.data
}

func (c *Complex) VariableType() Type {
	return ComplexType
}

func (c *Complex) String() string {
	return fmt.Sprintf("%v", c.data)
}

type Real struct {
	data float64
}

func NewReal(data float64) (*Real, error) {
	return &Real{data: data}, nil
}

func (r *Real) Data() float64 {
	return r.data
}

func (r *Real) VariableType() Type {
	return RealType
}

func (r *Real) String() string {
	return strconv.FormatFloat(r.data, 'f', -1, 64)
}

type Integer struct {
	data int64
}

func NewInteger(data int64) (*Integer, error) {
	return &Integer{data: data}, nil
}

func (i *Integer) Data() int64 {
	return i.data
}

func (i *Integer) VariableType() Type {
	return IntegerType
}

func (i *Integer) String() string {
	return strconv.Itoa(int(i.data))
}

func NewIntegerFromString(data string) (*Integer, error) {
	i, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return nil, ErrIntegerConversion
	}
	return NewInteger(i)
}

func NewRealFromString(data string) (*Real, error) {
	f, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return nil, ErrRealConversion
	}
	return NewReal(f)
}

func NewComplexFromString(data string) (*Complex, error) {
	c, err := strconv.ParseComplex(data, 64)
	if err != nil {
		return nil, ErrComplexConversion
	}

	return NewComplex(complex64(c))
}

func NewValue(fieldType string, data string) (Value, error) {
	var val Value
	var err error

	switch fieldType {
	case "integer":
		val, err = NewIntegerFromString(data)
		if err != nil {
			return nil, err
		}
	case "real":
		val, err = NewRealFromString(data)
		if err != nil {
			return nil, err
		}
	case "complex":
		val, err = NewComplexFromString(data)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrInvalidType
	}
	return val, nil
}
