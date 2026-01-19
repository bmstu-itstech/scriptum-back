package ports

import "github.com/bmstu-itstech/scriptum-back/internal/domain/value"

type PasswordHasher interface {
	Hash(password string) (value.HashedPassword, error)
	Verify(password string, hashed value.HashedPassword) bool
}
