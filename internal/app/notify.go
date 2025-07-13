package app

import (
	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/bmstu-itstech/scriptum-back/internal/service"

	"context"
)

type Notify struct {
	userRepo *service.MockUserRepo
	notifier scripts.Notifier
}

func NewNotifyUseCase(userRepo *service.MockUserRepo, notifier scripts.Notifier) (*Notify, error) {
	return &Notify{userRepo: userRepo, notifier: notifier}, nil
}

func (n *Notify) SendResult(ctx context.Context, userID scripts.UserID, res scripts.Result) error {
	user, err := n.userRepo.GetUser(userID)
	if err != nil {
		return err
	}
	return n.notifier.Notify(ctx, res, user.Email())
}
