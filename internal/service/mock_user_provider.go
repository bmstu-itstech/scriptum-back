package service

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type MockUserProvider struct {
	users map[int64]scripts.User
}

func NewMockUserProvider() (*MockUserProvider, error) {
	users := make(map[int64]scripts.User)
	u1, _ := scripts.NewUser(
		1,
		"Vasya Pupkin",
		"vasya@pupkin.com",
		true,
	)
	u2, _ := scripts.NewUser(
		2,
		"Petya Ivanov",
		"petya@ivanov.com",
		false,
	)
	users[1] = *u1
	users[2] = *u2
	return &MockUserProvider{users: users}, nil
}

func (m *MockUserProvider) User(_ context.Context, id scripts.UserID) (*scripts.User, error) {
	if user, ok := m.users[int64(id)]; ok {
		return &user, nil
	}
	return nil, scripts.ErrUserNotFound
}
