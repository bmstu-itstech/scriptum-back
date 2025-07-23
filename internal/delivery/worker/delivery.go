package delivery

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/bmstu-itstech/scriptum-back/internal/usecase"
)

const maxConcurrent = 10

type LaunchHandler struct {
	usecase    *usecase.JobRunUC
	subscriber message.Subscriber
	watLogger  watermill.LoggerAdapter
	sem        chan struct{}
}

func UnmarshalJob(data []byte) (scripts.Job, error) {
	type alias struct {
		JobID        scripts.JobID   `json:"job_id"`
		UserID       scripts.UserID  `json:"user_id"`
		In           scripts.Vector  `json:"in"`
		Command      string          `json:"command"`
		StartedAt    time.Time       `json:"started_at"`
		OutFields    []scripts.Field `json:"OutFields"`
		UserEmail    scripts.Email   `json:"user_email"`
		NeedToNotify bool            `json:"need_to_notify"`
	}

	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return scripts.Job{}, err
	}

	script, err := scripts.NewJob(a.JobID, a.UserID, a.In, a.Command, a.StartedAt, a.OutFields, a.UserEmail, a.NeedToNotify)
	return *script, err
}

func NewLaunchHandler(
	jobRunUC usecase.JobRunUC,
	subscriber message.Subscriber,
	watLogger watermill.LoggerAdapter,
) (*LaunchHandler, error) {
	return &LaunchHandler{
		usecase:    &jobRunUC,
		subscriber: subscriber,
		watLogger:  watLogger,
		sem:        make(chan struct{}, maxConcurrent),
	}, nil
}

func (l *LaunchHandler) Listen(ctx context.Context) {
	messages, err := l.subscriber.Subscribe(ctx, "script-start")
	if err != nil {
		l.watLogger.Error("Subscribe error", err, nil)
		return
	}

	l.sem <- struct{}{}

	go func() {
		defer func() {
			<-l.sem
			if r := recover(); r != nil {
				l.watLogger.Error("panic recovered in launch handler", nil, watermill.LogFields{"recover": r})
			}
		}()

		for {
			select {
			case <-ctx.Done():
				l.watLogger.Info("LaunchHandler stopped due to context cancel", nil)
				return

			case msg, ok := <-messages:
				if !ok {
					l.watLogger.Info("LaunchHandler channel closed", nil)
					return
				}

				go func(msg *message.Message) {
					req, err := UnmarshalJob(msg.Payload)
					if err != nil {
						l.watLogger.Error("Decode error", err, nil)
						msg.Nack()
						return
					}

					reqCtx := context.Background()
					job := usecase.JobToDTO(req)
					if err := l.usecase.ProcessLaunchRequest(reqCtx, job); err != nil {
						l.watLogger.Error("Process error", err, nil)
						msg.Nack()
						return
					}

					msg.Ack()
				}(msg)
			}
		}
	}()
}
