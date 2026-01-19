package value

import "github.com/bmstu-itstech/scriptum-back/internal/domain"

const MinPasswordLength = 8

type Password struct {
	s string
}

type HashedPassword []byte

func NewPassword(s string) (Password, error) {
	if s == "" {
		return Password{}, domain.NewInvalidInputError("empty-password", "expected not empty password")
	}

	if len(s) < MinPasswordLength {
		return Password{}, domain.NewInvalidInputError("too-short-password", "too short password")
	}

	return Password{s: s}, nil
}
