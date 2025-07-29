package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/google/uuid"
)

type WatermillDispatcher struct {
	publisher message.Publisher
}

func NewLauncher(publisher message.Publisher) (*WatermillDispatcher, error) {
	return &WatermillDispatcher{
		publisher: publisher,
	}, nil
}

func MarshalJob(job scripts.Job, needToNotify bool) ([]byte, error) {
	type Job struct {
		JobID        scripts.JobID       `json:"job_id"`
		OwnerID      scripts.UserID      `json:"owner_id"`
		ScriptID     scripts.ScriptID    `json:"script_id"`
		Input        []scripts.JSONValue `json:"in"`
		NeedToNotify bool                `json:"need_to_notify"`
		CreatedAt    time.Time           `json:"started_at"`
	}

	rawInputs := make([]scripts.JSONValue, 0, len(job.Input()))
	for _, v := range job.Input() {
		rawInputs = append(rawInputs, scripts.FromValue(v))
	}

	return json.Marshal(Job{
		JobID:        job.ID(),
		OwnerID:      job.OwnerID(),
		ScriptID:     job.ScriptID(),
		Input:        rawInputs,
		CreatedAt:    job.CreatedAt(),
		NeedToNotify: needToNotify,
	})
}

func (d *WatermillDispatcher) Start(ctx context.Context, request scripts.Job, needToNotify bool) error {
	payload, err := MarshalJob(request, needToNotify)
	if err == nil {
		msg := message.NewMessage(uuid.NewString(), payload)
		_ = d.publisher.Publish("script-start", msg)

	}

	return err
}
