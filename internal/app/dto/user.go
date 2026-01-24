package dto

import (
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
)

type User struct {
	ID        string
	Email     string
	Name      string
	Role      string
	CreatedAt time.Time
}

func UserToDTO(u *entity.User) User {
	return User{
		ID:        string(u.ID()),
		Email:     u.Email().String(),
		Name:      u.Name(),
		Role:      u.Role().String(),
		CreatedAt: u.CreatedAt(),
	}
}

func UsersToDTOs(us []*entity.User) []User {
	res := make([]User, len(us))
	for i, u := range us {
		res[i] = UserToDTO(u)
	}
	return res
}
