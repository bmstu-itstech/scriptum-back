package service

import "github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"

type EmailNotifier interface {
	scripts.Notifier
}

type EmailNotify struct{}

func NewEmailNotifier() *EmailNotify {
	return &EmailNotify{}
}

func (e *EmailNotify) Notify(r scripts.Result) error {

	return nil
}
