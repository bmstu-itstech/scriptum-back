package scripts

type Vector struct {
	values []Value
}

func (v *Vector) Values() []Value {
	return v.values
}

func NewVector(values []Value) (*Vector, error) {
	if len(values) == 0 {
		return nil, ErrEmptyVector
	}
	return &Vector{values: values}, nil
}

func (v *Vector) Add(value Value) {
	v.values = append(v.values, value)
}

type Value struct {
	VariableType Type
}

type Complex struct {
	Value
	data complex64
}

func (com *Complex) Data() complex64 {
	return com.data
}

func NewComplex(data complex64) (*Complex, error) {
	return &Complex{data: data}, nil
}

type Real struct {
	Value
	data float64
}

func (real *Real) Data() float64 {
	return real.data
}

func NewReal(data float64) (*Real, error) {
	return &Real{data: data}, nil
}

type Integer struct {
	Value
	data int64
}

func (inte *Integer) Data() int64 {
	return inte.data
}

func NewInteger(data int64) (*Integer, error) {
	return &Integer{data: data}, nil
}
