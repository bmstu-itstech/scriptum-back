package scripts

type TypeValue int

const (
	INTEGER TypeValue = iota
	REAL
	COMPLEX
)

type Type struct {
	Value TypeValue
}
