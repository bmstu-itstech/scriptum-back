package watermill

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/ThreeDotsLabs/watermill/message"
)

type JobSubscriber struct {
	s message.Subscriber
	l *slog.Logger
}

func NewJobSubscriber(s message.Subscriber, l *slog.Logger) JobSubscriber {
	return JobSubscriber{s, l}
}

func (s JobSubscriber) Listen(ctx context.Context, callback func(ctx context.Context, jobID string) error) error {
	l := s.l.With(slog.String("op", "watermill.JobSubscriber.Listen"))

	msgCh, err := s.s.Subscribe(ctx, topicRunJob)
	if err != nil {
		l.ErrorContext(
			ctx,
			"failed to subscribe to topic run job subscriber",
			slog.String("topic", topicRunJob),
			slog.String("error", err.Error()),
		)
		return err
	}

	l.InfoContext(ctx, "successfully subscribed to topic run job subscriber, listening")
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case msg, ok := <-msgCh:
			if !ok {
				// Закрытие канала
				l.InfoContext(ctx, "channel closed")
				return nil
			}
			go func() { s.handle(ctx, msg, callback) }()
		}
	}
}

func (s JobSubscriber) handle(
	ctx context.Context,
	msg *message.Message,
	callback func(ctx context.Context, jobID string) error,
) {
	l := s.l.With(
		slog.String("op", "watermill.JobSubscriber.handle"),
		slog.String("message", msg.UUID),
	)

	var pl payload
	if err := json.Unmarshal(msg.Payload, &pl); err != nil {
		l.ErrorContext(ctx, "failed to unmarshal payload", slog.String("error", err.Error()))
		return
	}

	l = l.With(slog.String("jobID", pl.JobID))
	if err := callback(context.Background(), pl.JobID); err != nil {
		l.ErrorContext(ctx, "failed to handle payload", slog.String("error", err.Error()))
		return
	}
	l.InfoContext(ctx, "handled job")
}
