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

func (u *User) BlueprintVisibility() value.Visibility {
	switch u.role {
	case value.RoleAdmin:
		return value.VisibilityPublic
	}
	return value.VisibilityPrivate
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

func (u *User) Role() value.Role {
	return u.role
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
