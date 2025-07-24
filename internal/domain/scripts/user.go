package scripts

import "fmt"

type Name string
type UserID uint32
type Email string

type User struct {
	userID  UserID
	name    Name
	email   Email
	isAdmin bool
}

func (u *User) UserID() UserID {
	return u.userID
}

func (u *User) Name() Name {
	return u.name
}

func (u *User) Email() Email {
	return u.email
}

func (u *User) IsAdmin() bool {
	return u.isAdmin
}

func NewUser(userID UserID, fullName Name, email Email, isAdmin bool) (*User, error) {
	if fullName == "" {
		return nil, fmt.Errorf("%w: invalid User: expected not empty name", ErrInvalidInput)
	}

	if email == "" {
		return nil, fmt.Errorf("%w: invalid User: expected not empty email", ErrInvalidInput)
	}

	return &User{
		userID:  userID,
		name:    fullName,
		email:   email,
		isAdmin: isAdmin,
	}, nil
}
