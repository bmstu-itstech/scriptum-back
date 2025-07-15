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
	return &Complex{
		data: data,
	}, nil
}

func (c *Complex) VariableType() Type {
	return ComplexType
}

func (c *Complex) String() string {
	return fmt.Sprintf("Complex(%v)", c.data)
}

func (c *Complex) Data() complex64 {
	return c.data
}

type Real struct {
	data float64
}

func (r *Real) VariableType() Type {
	return RealType
}

func (r *Real) String() string {
	return fmt.Sprintf("Real(%f)", r.data)
}

func (r *Real) Data() float64 {
	return r.data
}

func NewReal(data float64) (*Real, error) {
	return &Real{
		data: data,
	}, nil
}

type Integer struct {
	data int64
}

func (i *Integer) VariableType() Type {
	return IntegerType
}

func (i *Integer) String() string {
	return fmt.Sprintf("Integer(%d)", i.data)
}

func (i *Integer) Data() int64 {
	return i.data
}

func NewInteger(data int64) (*Integer, error) {
	return &Integer{
		data: data,
	}, nil
}

func NewIntegerString(data string) (*Integer, error) {
	i, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return nil, ErrIntegerConversion
	}
	return NewInteger(i)
}

func NewRealString(data string) (*Real, error) {
	f, err := strconv.ParseFloat(data, 64)
	if err != nil {
		return nil, ErrRealConversion
	}
	return NewReal(f)
}

func NewComplexString(data string) (*Complex, error) {
	var r, i float64

	n, err := fmt.Sscanf(data, "%f+%fi", &r, &i)
	if err != nil || n != 2 {
		n, err = fmt.Sscanf(data, "%f-%fi", &r, &i)
		if err != nil || n != 2 {
			return nil, ErrComplexConversion
		}
		i = -i
	}

	c := complex(float32(r), float32(i))
	return NewComplex(c)
}

func NewValue(fieldType string, data string) (Value, error) {
	var val Value
	var err error

	switch fieldType {
	case "integer":
		val, err = NewIntegerString(data)
		if err != nil {
			return nil, err
		}
	case "real":
		val, err = NewRealString(data)
		if err != nil {
			return nil, err
		}
	case "complex":
		val, err = NewComplexString(data)
		if err != nil {
			return nil, err
		}
	}
	return val, nil
}
