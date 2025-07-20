package delivery

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	worker "github.com/bmstu-itstech/scriptum-back/internal/usecase/worker"
)

type LaunchHandler struct {
	usecase    *worker.JobLaunchUC
	subscriber message.Subscriber
	watLogger  watermill.LoggerAdapter
}

func NewLaunchHandler(
	jobR scripts.JobRepository,
	launcher scripts.Launcher,
	notifier scripts.Notifier,
	subscriber message.Subscriber,
	watLogger watermill.LoggerAdapter,
) (*LaunchHandler, error) {
	usecase, err := worker.NewJobLaunchUC(jobR, launcher, notifier)
	if err != nil {
		return nil, err
	}
	return &LaunchHandler{
		usecase:    usecase,
		subscriber: subscriber,
		watLogger:  watLogger,
	}, nil
}

func (l *LaunchHandler) Listen(ctx context.Context) {
	messages, err := l.subscriber.Subscribe(ctx, "script-start")
	if err != nil {
		l.watLogger.Error("Subscribe error", err, nil)
		return
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				l.watLogger.Error("panic recovered in email notifier", nil, watermill.LogFields{"recover": r})
			}
		}()

		for {
			select {
			case <-ctx.Done():
				l.watLogger.Info("EmailNotifier stopped due to context cancel", nil)
				return
			case msg, ok := <-messages:
				if !ok {
					l.watLogger.Info("EmailNotifier channel closed", nil)
					return
				}
				var req scripts.LaunchRequest

				if err := req.UnmarshalJSON(msg.Payload); err != nil {
					l.watLogger.Error("Decode error", err, nil)
					msg.Nack()
					continue
				}

				ctx := context.Background()

				if err := l.usecase.ProcessLaunchRequest(ctx, req); err != nil {
					l.watLogger.Error("Process error", err, nil)
					msg.Nack()
					continue
				}

				msg.Ack()
			}
		}
	}()
}
