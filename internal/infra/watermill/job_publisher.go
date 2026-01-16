package watermill

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
)

type JobPublisher struct {
	p message.Publisher
}

func NewJobPublisher(p message.Publisher) JobPublisher {
	return JobPublisher{p}
}

func (p JobPublisher) PublishJob(_ context.Context, job *entity.Job) error {
	pl := payload{
		JobID: string(job.ID()),
	}
	msg, err := json.Marshal(pl)
	if err != nil {
		return fmt.Errorf("failed to marshall job: %w", err)
	}
	wMsg := message.NewMessage(
		watermill.NewShortUUID(),
		msg,
	)
	return p.p.Publish(topicRunJob, wMsg)
}
