package scripts

import "fmt"

type Name = string
type UserID = uint32
type Email = string

type User struct {
	userID   UserID
	fullName Name
	email    Email
	isAdmin  bool
}

func (u *User) UserID() UserID {
	return u.userID
}

func (u *User) FullName() Name {
	return u.fullName
}

func (u *User) Email() Email {
	return u.email
}

func (u *User) IsAdmin() bool {
	return u.isAdmin
}

func NewUser(userID UserID, fullName Name, email Email, isAdmin bool) (*User, error) {
	if fullName == "" {
		return nil, fmt.Errorf("fullName: expected non-empty string, got empty string: %w", ErrUserInvalid)
	}
	if email == "" {
		return nil, fmt.Errorf("email: expected non-empty string, got empty string: %w", ErrUserInvalid)
	}

	return &User{
		userID:   userID,
		fullName: fullName,
		email:    email,
		isAdmin:  isAdmin,
	}, nil
}
