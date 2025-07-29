package scripts

import (
	"fmt"
	"strconv"
)

type Value interface {
	Type() ValueType
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

func (c *Complex) Type() ValueType {
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

func (r *Real) Type() ValueType {
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

func (i *Integer) Type() ValueType {
	return IntegerType
}

func (i *Integer) String() string {
	return strconv.Itoa(int(i.data))
}

func NewIntegerFromString(data string) (*Integer, error) {
	i, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid Integer: expected valid integer, got %s", ErrInvalidInput, err)
	}
	return NewInteger(i)
}

func NewRealFromString(data string) (*Real, error) {
	f, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid Real: expected valid real, got %s", ErrInvalidInput, err)
	}
	return NewReal(f)
}

func NewComplexFromString(data string) (*Complex, error) {
	c, err := strconv.ParseComplex(data, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid Complex: expected valid complex, got %s", ErrInvalidInput, err)
	}

	return NewComplex(complex64(c))
}

func NewValue(typ string, data string) (Value, error) {
	switch typ {
	case "integer":
		return NewIntegerFromString(data)

	case "real":
		return NewRealFromString(data)

	case "complex":
		return NewComplexFromString(data)
	}

	return nil, fmt.Errorf(
		"%w: invalid Value: expected type one of ['integer', 'real', 'complex'], got %s",
		ErrInvalidInput, typ,
	)
}
