package service

import (
	"context"
	"encoding/json"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/google/uuid"
)

type Launcher struct {
	publisher message.Publisher
}

type ScriptFinishedEvent struct {
	result    scripts.Result
	userEmail scripts.Email
}

func NewLauncher(interpreter string, publisher message.Publisher, flags ...string) (*Launcher, error) {
	return &Launcher{
		publisher: publisher,
	}, nil
}

func (p *Launcher) Launch(ctx context.Context, job scripts.Job, scriptFields []scripts.Field, userEmail scripts.Email, needToNotify bool) error {
	// создать сообшение, результат не пришел

	// пишем в очередь сообщений

	request := scripts.NewLaunchRequest(job, scriptFields, userEmail, needToNotify)

	payload, err := json.Marshal(request)
	if err == nil {
		msg := message.NewMessage(uuid.NewString(), payload)
		_ = p.publisher.Publish("script-finished", msg)

	}

	return nil
}
