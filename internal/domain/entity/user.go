package entity

import (
	"errors"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type User struct {
	id        value.UserID
	email     value.Email
	passhash  []byte
	name      string
	role      value.Role
	createdAt time.Time
}

func NewUser(name string, email value.Email, passhash []byte, role value.Role) (*User, error) {
	if name == "" {
		return nil, domain.NewInvalidInputError("user-empty-name", "expected not empty name")
	}

	if email.IsZero() {
		return nil, domain.NewInvalidInputError("user-empty-email", "expected not empty email")
	}

	if len(passhash) == 0 {
		return nil, errors.New("empty passhash")
	}

	if role.IsZero() {
		return nil, errors.New("zero role")
	}

	id := value.NewUserID()
	return &User{
		id:        id,
		email:     email,
		passhash:  passhash,
		name:      name,
		role:      role,
		createdAt: time.Now(),
	}, nil
}

func (u *User) CanCreateBlueprintWithVisibility(v value.Visibility) bool {
	switch u.role {
	case value.RoleAdmin:
		return true
	case value.RoleUser:
		return v == value.VisibilityPrivate
	}
	return false
}

func (u *User) CanSee(uid value.UserID) bool {
	if u.role == value.RoleAdmin {
		return true
	}
	return u.id == uid
}

func (u *User) SetEmail(email value.Email) error {
	if email.IsZero() {
		return domain.NewInvalidInputError("user-empty-email", "expected not empty email")
	}
	u.email = email
	return nil
}

func (u *User) SetPassword(password []byte) error {
	if len(password) == 0 {
		return domain.NewInvalidInputError("user-empty-password", "expected not empty password")
	}
	u.passhash = password
	return nil
}

func (u *User) SetName(name string) error {
	if name == "" {
		return domain.NewInvalidInputError("user-empty-name", "expected not empty name")
	}
	u.name = name
	return nil
}

func (u *User) SetRole(role value.Role) error {
	if role.IsZero() {
		return domain.NewInvalidInputError("user-empty-role", "expected not empty role")
	}
	u.role = role
	return nil
}

func (u *User) ID() value.UserID {
	return u.id
}

func (u *User) Name() string {
	return u.name
}

func (u *User) Email() value.Email {
	return u.email
}

func (u *User) PasswordHash() []byte {
	return u.passhash
}

func (u *User) Role() value.Role {
	return u.role
}

func (u *User) CreatedAt() time.Time {
	return u.createdAt
}

func RestoreUser(
	id value.UserID,
	email value.Email,
	passhash []byte,
	name string,
	role value.Role,
	createdAt time.Time,
) (*User, error) {
	if id == "" {
		return nil, errors.New("empty user id")
	}

	if email.IsZero() {
		return nil, errors.New("empty email")
	}

	if name == "" {
		return nil, errors.New("empty name")
	}

	if role.IsZero() {
		return nil, errors.New("empty role")
	}

	if createdAt.IsZero() {
		return nil, errors.New("empty createdAt")
	}

	return &User{
		id:        id,
		email:     email,
		passhash:  passhash,
		name:      name,
		role:      role,
		createdAt: createdAt,
	}, nil
}
