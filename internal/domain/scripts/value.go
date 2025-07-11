package scripts

type Vector struct {
	Values []Value
}

type Value struct {
	VariableType Type
}

type Complex struct {
	Value
	Data complex64
}

type Real struct {
	Value
	Data float32
}

type Integer struct {
	Value
	Data int
}
