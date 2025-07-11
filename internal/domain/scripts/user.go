package scripts

type Name = string
type UserID = uint32
type Email = string

type User struct {
	UserID   UserID
	fullName Name
	email    Email
	isAdmin  bool
}

func (u User) FullName() Name {
	return u.fullName
}

func (u User) Email() Email {
	return u.email
}

func (u User) IsAdmin() bool {
	return u.isAdmin
}
