package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type MockUserRepo struct {
	sync.RWMutex
	m         map[scripts.UserID]scripts.User
	idCounter scripts.UserID
}

func NewMockUserRepository() (*MockUserRepo, error) {
	return &MockUserRepo{
		m: make(map[scripts.UserID]scripts.User),
	}, nil
}

func (r *MockUserRepo) User(ctx context.Context, id scripts.UserID) (*scripts.User, error) {
	r.RLock()
	defer r.RUnlock()

	user, ok := r.m[id]
	if !ok {
		return nil, fmt.Errorf("user not found: %d", id)
	}
	return &user, nil
}

func (r *MockUserRepo) Users(ctx context.Context) ([]scripts.User, error) {
	r.RLock()
	defer r.RUnlock()

	users := make([]scripts.User, 0, len(r.m))
	for _, u := range r.m {
		users = append(users, u)
	}
	return users, nil
}

func (r *MockUserRepo) StoreUser(ctx context.Context, user scripts.User) (scripts.UserID, error) {
	r.Lock()
	defer r.Unlock()

	r.idCounter++
	r.m[user.UserID()] = user
	return user.UserID(), nil
}

func (r *MockUserRepo) DeleteUser(ctx context.Context, userID scripts.UserID) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[userID]; !ok {
		return fmt.Errorf("user not found: %d", userID)
	}
	delete(r.m, userID)
	return nil
}

func (r *MockUserRepo) UpdateUser(ctx context.Context, user scripts.User) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[user.UserID()]; !ok {
		return fmt.Errorf("user not found: %d", user.UserID())
	}
	r.m[user.UserID()] = user
	return nil
}
