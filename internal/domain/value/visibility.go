package value

type Visibility struct {
	s string
}

var (
	VisibilityPublic  = Visibility{s: "public"}
	VisibilityPrivate = Visibility{s: "private"}
)

func (v Visibility) String() string {
	return v.s
}

func (v Visibility) IsZero() bool {
	return v.s == ""
}
