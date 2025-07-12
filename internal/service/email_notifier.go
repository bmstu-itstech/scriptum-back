package service

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type EmailNotifier struct{}

func NewEmailNotifier() (*EmailNotifier, error) {
	return &EmailNotifier{}, nil
}

func (e *EmailNotifier) Notify(_ context.Context, r scripts.Result) error {

	return nil
}
