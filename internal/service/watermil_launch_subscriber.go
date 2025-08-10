package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/bmstu-itstech/scriptum-back/internal/app"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type LaunchSubscriber struct {
	subscriber message.Subscriber
	watLogger  watermill.LoggerAdapter
}

func UnmarshalJob(data []byte) (*app.JobDTO, error) {
	var jsonJob WJob
	if err := json.Unmarshal(data, &jsonJob); err != nil {
		return nil, err
	}

	input := make([]scripts.Value, len(jsonJob.Input))
	for i, v := range jsonJob.Input {
		val, err := v.toValue()
		if err != nil {
			return nil, err
		}
		input[i] = val
	}

	inputVal, err := app.ValuesToDTO(input)
	if err != nil {
		return nil, err
	}

	exp := make([]scripts.Field, len(jsonJob.Expected))
	for i, v := range jsonJob.Expected {
		val, err := scripts.NewValueType(v.Type)
		if err != nil {
			return nil, err
		}
		f, err := scripts.NewField(*val, v.Name, v.Desc, v.Unit)
		if err != nil {
			return nil, err
		}
		exp[i] = *f
	}

	expected, err := app.FieldsToDTO(exp)
	if err != nil {
		return nil, err
	}

	job := &app.JobDTO{
		JobID:        int64(jsonJob.JobID),
		OwnerID:      int64(jsonJob.OwnerID),
		ScriptID:     int64(jsonJob.ScriptID),
		Url:          jsonJob.URL,
		Input:        inputVal,
		Expected:     expected,
		State:        scripts.JobPending.String(),
		CreatedAt:    jsonJob.CreatedAt,
		FinishedAt:   nil,
		NeedToNotify: jsonJob.NeedToNotify,
	}

	return job, nil
}

func NewLaunchHandler(
	subscriber message.Subscriber,
	watLogger watermill.LoggerAdapter,
) (*LaunchSubscriber, error) {
	return &LaunchSubscriber{
		subscriber: subscriber,
		watLogger:  watLogger,
	}, nil
}

func (l *LaunchSubscriber) Listen(ctx context.Context, callback func(context.Context, app.JobDTO) error) error {
	messages, err := l.subscriber.Subscribe(ctx, "script-start")
	if err != nil {
		l.watLogger.Error("Subscribe error", err, nil)
		return err
	}

	go func() {
		defer func() {
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
						fmt.Println(err)
						l.watLogger.Error("Decode error", err, nil)
						msg.Nack()
						return
					}
					msg.Ack()

					if err := callback(ctx, *req); err != nil {
						l.watLogger.Error("Callback error", err, nil)
						msg.Nack()
						return
					}
				}(msg)
			}
		}
	}()
	return nil
}
