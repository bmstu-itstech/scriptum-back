package value

import (
	"fmt"

	"github.com/bmstu-itstech/scriptum-back/internal/domain"
)

type Role struct {
	s string
}

var RoleUser = Role{"user"}
var RoleAdmin = Role{"admin"}

func RoleFromString(s string) (Role, error) {
	switch s {
	case "user":
		return RoleUser, nil
	case "admin":
		return RoleAdmin, nil
	}
	return Role{}, domain.NewInvalidInputError(
		"role-invalid",
		fmt.Sprintf("invalid value type: expected one of ['user' 'admin'], got %s", s),
	)
}

func (r Role) String() string {
	return r.s
}

func (r Role) IsZero() bool {
	return r == Role{}
}
