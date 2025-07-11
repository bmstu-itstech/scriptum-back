package service

import (
	"fmt"
	"sync"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type MockUserRepo struct {
	sync.RWMutex
	m map[scripts.UserID]scripts.User
}

func NewMockUserRepository() (*MockUserRepo, error) {
	return &MockUserRepo{
		m: make(map[scripts.UserID]scripts.User),
	}, nil
}

func (r *MockUserRepo) GetUser(userId scripts.UserID) (scripts.User, error) {
	r.RLock()
	defer r.RUnlock()

	user, ok := r.m[userId]
	if !ok {
		return scripts.User{}, fmt.Errorf("user not found: %d", userId)
	}
	return user, nil
}
