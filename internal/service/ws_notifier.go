package service

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type WsNotifier struct{}

func NewWsNotifier() (*WsNotifier, error) {
	return &WsNotifier{}, nil
}

func (w *WsNotifier) Notify(_ context.Context, r scripts.Result, email scripts.Email) error {

	return nil
}
