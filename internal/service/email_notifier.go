package service

import "github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"

type EmailNotifier struct{}

func NewEmailNotifier() (*EmailNotifier, error) {
	return &EmailNotifier{}, nil
}

func (e *EmailNotifier) Notify(r scripts.Result) error {

	return nil
}
