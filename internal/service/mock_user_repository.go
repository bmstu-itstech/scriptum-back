package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type MockUserRepository interface {
	GetUser(userId scripts.UserID) (scripts.User, error)
}

type MockUserRepo struct {
	context context.Context
	sync.RWMutex
	m map[scripts.UserID]scripts.User
}

func NewMockUserRepository() *MockUserRepo {
	return &MockUserRepo{
		context: context.Background(),
		m:       make(map[scripts.UserID]scripts.User),
	}
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
