package mock

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type IsAdminProvider struct {
	admins map[value.UserID]bool
}

func NewIsAdminProvider() *IsAdminProvider {
	return &IsAdminProvider{
		admins: make(map[value.UserID]bool),
	}
}

func (i *IsAdminProvider) IsAdmin(_ context.Context, uid value.UserID) (bool, error) {
	return i.admins[uid], nil
}
