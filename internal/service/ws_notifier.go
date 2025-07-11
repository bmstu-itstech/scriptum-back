package service

import "github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"

type WsNotifier struct{}

func NewWsNotifier() *WsNotifier {
	return &WsNotifier{}
}

func (w *WsNotifier) Notify(r scripts.Result) error {

	return nil
}
