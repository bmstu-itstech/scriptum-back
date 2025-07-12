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

func (v *Vector) Len() int {
	return len(v.values)
}
