package service

import (
	"context"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/google/uuid"
)

type Dispatcher struct {
	publisher message.Publisher
}

func NewLauncher(publisher message.Publisher) (*Dispatcher, error) {
	return &Dispatcher{
		publisher: publisher,
	}, nil
}

func (d *Dispatcher) Launch(ctx context.Context, request scripts.LaunchRequest) error {
	payload, err := request.MarshalJSON()
	if err == nil {
		msg := message.NewMessage(uuid.NewString(), payload)
		_ = d.publisher.Publish("script-start", msg)

	}

	return err
}
