package scripts

import "fmt"

type Value interface {
	VariableType() Type
	String() string
}

type Complex struct {
	variableType Type
	data         complex64
}

func NewComplex(data complex64) (*Complex, error) {
	t, _ := NewType("complex")
	return &Complex{variableType: *t, data: data}, nil
}

func (c *Complex) VariableType() Type {
	return c.variableType
}

func (c *Complex) Data() complex64 {
	return c.data
}

func (c *Complex) String() string {
	return fmt.Sprintf("Complex(%v)", c.data)
}

type Real struct {
	variableType Type
	data         float64
}

func NewReal(data float64) (*Real, error) {
	t, _ := NewType("real")
	return &Real{variableType: *t, data: data}, nil
}

func (r *Real) VariableType() Type {
	return r.variableType
}

func (r *Real) Data() float64 {
	return r.data
}

func (r *Real) String() string {
	return fmt.Sprintf("Real(%f)", r.data)
}

type Integer struct {
	variableType Type
	data         int64
}

func NewInteger(data int64) (*Integer, error) {
	t, _ := NewType("integer")
	return &Integer{variableType: *t, data: data}, nil
}

func (i *Integer) VariableType() Type {
	return i.variableType
}

func (i *Integer) Data() int64 {
	return i.data
}

func (i *Integer) String() string {
	return fmt.Sprintf("Integer(%d)", i.data)
}
