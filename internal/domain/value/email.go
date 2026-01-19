package value

import (
	"net/mail"

	"github.com/bmstu-itstech/scriptum-back/internal/domain"
)

type Email struct {
	s string
}

func EmailFromString(s string) (Email, error) {
	_, err := mail.ParseAddress(s)
	if err != nil {
		return Email{}, domain.NewInvalidInputError("invalid-email", "email is invalid")
	}
	return Email{s}, nil
}

func MustEmailFromString(s string) Email {
	email, err := EmailFromString(s)
	if err != nil {
		panic(err)
	}
	return email
}

func (e Email) String() string {
	return e.s
}

func (e Email) IsZero() bool {
	return e.s == ""
}
