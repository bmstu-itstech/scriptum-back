package value

import "strings"

type Input struct {
	b strings.Builder
}

func NewEmptyInput() Input {
	return Input{}
}

func (i Input) With(v Value) Input {
	i.b.WriteString(v.String())
	i.b.WriteRune('\n')
	return i
}

func (i Input) String() string {
	return i.b.String()
}
