package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/google/uuid"
)

type LaunchPublisher struct {
	publisher message.Publisher
}

func NewLauncher(publisher message.Publisher) (*LaunchPublisher, error) {
	return &LaunchPublisher{
		publisher: publisher,
	}, nil
}

func MarshalJob(job scripts.Job, needToNotify bool) ([]byte, error) {
	type Job struct {
		JobID    scripts.JobID    `json:"job_id"`
		OwnerID  scripts.UserID   `json:"owner_id"`
		ScriptID scripts.ScriptID `json:"script_id"`
		Input    []JSONValue      `json:"in"`
		Expected []scripts.Field  `json:"exp"`
		URL      string           `json:"url"`

		NeedToNotify bool      `json:"need_to_notify"`
		CreatedAt    time.Time `json:"started_at"`
	}

	rawInputs := make([]JSONValue, 0, len(job.Input()))
	for _, v := range job.Input() {
		rawInputs = append(rawInputs, fromValue(v))
	}

	return json.Marshal(Job{
		JobID:        job.ID(),
		OwnerID:      job.OwnerID(),
		ScriptID:     job.ScriptID(),
		Input:        rawInputs,
		Expected:     job.Expected(),
		URL:          job.URL(),
		CreatedAt:    job.CreatedAt(),
		NeedToNotify: needToNotify,
	})
}

func (d *LaunchPublisher) Start(ctx context.Context, request *scripts.Job, needToNotify bool) error {
	payload, err := MarshalJob(*request, needToNotify)
	if err == nil {
		msg := message.NewMessage(uuid.NewString(), payload)
		err = d.publisher.Publish("script-start", msg)
	}

	return err
}

type JSONValue struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (jv *JSONValue) toValue() (scripts.Value, error) {
	return scripts.NewValue(jv.Type, jv.Value)
}

func fromValue(v scripts.Value) JSONValue {
	return JSONValue{
		Type:  v.Type().String(),
		Value: v.String(),
	}
}
