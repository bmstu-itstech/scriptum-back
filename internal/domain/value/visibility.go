package value

import (
	"fmt"

	"github.com/bmstu-itstech/scriptum-back/internal/domain"
)

type Visibility struct {
	s string
}

var (
	VisibilityPublic  = Visibility{s: "public"}
	VisibilityPrivate = Visibility{s: "private"}
)

func VisibilityFromString(s string) (Visibility, error) {
	switch s {
	case "public":
		return VisibilityPublic, nil
	case "private":
		return VisibilityPrivate, nil
	}
	return Visibility{}, domain.NewInvalidInputError(
		"visibility-invalid",
		fmt.Sprintf("invalid visibility: expected one of ['public', 'private'], got '%s'", s),
	)
}

func (v Visibility) String() string {
	return v.s
}

func (v Visibility) IsZero() bool {
	return v.s == ""
}
