package scripts

import "fmt"

type Value struct {
	variableType Type
}

func (v *Value) VariableType() Type {
	return v.variableType
}

type Complex struct {
	Value
	data complex64
}

func (c *Complex) Data() complex64 {
	return c.data
}

func NewComplex(data complex64) (*Complex, error) {
	return &Complex{data: data}, nil
}

type Real struct {
	Value
	data float64
}

func (r *Real) Data() float64 {
	return r.data
}

func NewReal(data float64) (*Real, error) {
	return &Real{data: data}, nil
}

type Integer struct {
	Value
	data int64
}

func (i *Integer) Data() int64 {
	return i.data
}

func NewInteger(data int64) (*Integer, error) {
	return &Integer{data: data}, nil
}

func (c *Complex) String() string {
	return fmt.Sprintf("Complex(%v)", c.data)
}

func (r *Real) String() string {
	return fmt.Sprintf("Real(%f)", r.data)
}

func (i *Integer) String() string {
	return fmt.Sprintf("Integer(%d)", i.data)
}
