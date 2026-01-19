package bcrypt

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type PasswordHasher struct {
	cost int
}

func NewPasswordHasher(cost int) *PasswordHasher {
	return &PasswordHasher{cost}
}

func (p *PasswordHasher) Hash(password string) (value.HashedPassword, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), p.cost)
	if err != nil {
		return value.HashedPassword{}, err
	}
	return hash, nil
}

func (p *PasswordHasher) Verify(password string, hashed value.HashedPassword) bool {
	err := bcrypt.CompareHashAndPassword(hashed, []byte(password))
	return err == nil
}
