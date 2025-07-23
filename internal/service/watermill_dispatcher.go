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

func MarshalJob(curJob scripts.Job) ([]byte, error) {
	type job struct {
		JobID        scripts.JobID   `json:"job_id"`
		UserID       scripts.UserID  `json:"user_id"`
		In           scripts.Vector  `json:"in"`
		Command      string          `json:"command"`
		StartedAt    time.Time       `json:"started_at"`
		OutFields    []scripts.Field `json:"out_fields"`
		UserEmail    scripts.Email   `json:"user_email"`
		NeedToNotify bool            `json:"need_to_notify"`
	}

	return json.Marshal(job{
		JobID:        curJob.JobID(),
		UserID:       curJob.UserID(),
		In:           curJob.In(),
		Command:      curJob.Command(),
		StartedAt:    curJob.StartedAt(),
		OutFields:    curJob.OutFields(),
		UserEmail:    curJob.UserEmail(),
		NeedToNotify: curJob.NeedToNotify(),
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
