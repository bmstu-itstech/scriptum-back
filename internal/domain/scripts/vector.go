package scripts

type Vector struct {
	values []Value
}

func (v *Vector) Values() []Value {
	return v.values
}

func NewVector(values []Value) (*Vector, error) {
	if len(values) == 0 {
		return nil, ErrVectorEmpty
	}
	return &Vector{values: values}, nil
}

func (v *Vector) Add(value Value) {
	v.values = append(v.values, value)
}

func (v *Vector) Len() int {
	return len(v.values)
}

func (v *Vector) Get() []string {
	values := []string{}
	for _, v := range v.Values() {
		switch v.VariableType() {
		case ComplexType:
			values = append(values, v.String())
		case RealType:
			values = append(values, v.String())
		case IntegerType:
			values = append(values, v.String())
		default:
			values = append(values, "unknown value")
		}
	}
	return values
}
