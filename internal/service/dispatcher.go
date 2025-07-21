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

func MarshalJob(job scripts.Job) ([]byte, error) {
	type alias struct {
		JobID        scripts.JobID   `json:"job_id"`
		UserID       scripts.UserID  `json:"user_id"`
		In           scripts.Vector  `json:"in"`
		Command      string          `json:"command"`
		StartedAt    time.Time       `json:"started_at"`
		ScriptFields []scripts.Field `json:"script_fields"`
		UserEmail    scripts.Email   `json:"user_email"`
		NeedToNotify bool            `json:"need_to_notify"`
	}

	return json.Marshal(alias{
		JobID:        job.JobID(),
		UserID:       job.UserID(),
		In:           job.In(),
		Command:      job.Command(),
		StartedAt:    job.StartedAt(),
		ScriptFields: job.ScriptFields(),
		UserEmail:    job.UserEmail(),
		NeedToNotify: job.NeedToNotify(),
	})
}

func (d *WatermillDispatcher) Launch(ctx context.Context, request scripts.Job) error {
	payload, err := MarshalJob(request)
	if err == nil {
		msg := message.NewMessage(uuid.NewString(), payload)
		_ = d.publisher.Publish("script-start", msg)

	}

	return err
}
