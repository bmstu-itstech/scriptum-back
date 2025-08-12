package service

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/google/uuid"
)

type WJob struct {
	JobID    scripts.JobID    `json:"job_id"`
	OwnerID  scripts.UserID   `json:"owner_id"`
	ScriptID scripts.ScriptID `json:"script_id"`
	Input    []JSONValue      `json:"in"`
	Expected []JSONField      `json:"exp"`
	URL      string           `json:"url"`

	NeedToNotify bool      `json:"need_to_notify"`
	CreatedAt    time.Time `json:"started_at"`
}

type LaunchPublisher struct {
	publisher message.Publisher
}

func NewLauncher(publisher message.Publisher) (*LaunchPublisher, error) {
	return &LaunchPublisher{
		publisher: publisher,
	}, nil
}

type JSONField struct {
	Type string `json:"type"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Unit string `json:"unit"`
}

func fromField(f scripts.Field) JSONField {
	return JSONField{
		Type: f.ValueType().String(),
		Name: f.Name(),
		Desc: f.Description(),
		Unit: f.Unit(),
	}
}

func marshalJob(job scripts.Job, needToNotify bool) ([]byte, error) {
	rawInputs := make([]JSONValue, 0, len(job.Input()))
	for _, v := range job.Input() {
		rawInputs = append(rawInputs, fromValue(v))
	}

	rawExp := make([]JSONField, 0, len(job.Expected()))
	for _, f := range job.Expected() {
		rawExp = append(rawExp, fromField(f))
	}

	return json.Marshal(WJob{
		JobID:        job.ID(),
		OwnerID:      job.OwnerID(),
		ScriptID:     job.ScriptID(),
		Input:        rawInputs,
		Expected:     rawExp,
		URL:          job.URL(),
		CreatedAt:    job.CreatedAt(),
		NeedToNotify: needToNotify,
	})
}

func (d *LaunchPublisher) Start(ctx context.Context, request *scripts.Job, needToNotify bool) error {
	payload, err := marshalJob(*request, needToNotify)
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
