package app

import (
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type FieldDTO struct {
	Type string
	Name string
	Desc string
	Unit string
}

type ScriptDTO struct {
	ID         int32
	OwnerID    int64
	Name       string
	Desc       string
	Input      []FieldDTO
	Output     []FieldDTO
	URL        string
	Visibility string
	CreatedAt  time.Time
}

type ValueDTO struct {
	Type string
	Data string
}

type JobDTO struct {
	JobID        int64
	OwnerID      int64
	ScriptID     int64
	Input        []ValueDTO
	State        string
	CreatedAt    time.Time
	FinishedAt   *time.Time
	NeedToNotify bool
}

type FileDTO struct {
	Name    string
	Content []byte
}

type ScriptCreateDTO struct {
	OwnerID           int64
	ScriptName        string
	ScriptDescription string
	File              FileDTO
	InFields          []FieldDTO
	OutFields         []FieldDTO
}

type ScriptRunDTO struct {
	ScriptID     uint32
	InParams     []ValueDTO
	NeedToNotify bool
}

func DTOToFields(_ []FieldDTO) ([]scripts.Field, error) {
	return nil, nil
}

func ScriptToDTO(_ scripts.Script) ScriptDTO {
	return ScriptDTO{}
}

func DTOToJob(j JobDTO) (*scripts.Job, error) {
	values, err := DTOToValues(j.Input)
	if err != nil {
		return nil, err
	}

	job, err := scripts.RestoreJob(
		j.JobID,
		j.OwnerID,
		j.ScriptID,
		j.State,
		values,
		nil,
		j.CreatedAt,
		j.FinishedAt,
	)

	return job, err
}

func ValuesToDTO(values []scripts.Value) ([]ValueDTO, error) {
	jobValues := make([]ValueDTO, len(values))
	for i, v := range values {
		val := ValueDTO{
			v.Type().String(),
			v.String(),
		}
		jobValues[i] = val
	}
	return jobValues, nil
}

func DTOToValues(values []ValueDTO) ([]scripts.Value, error) {
	jobValues := make([]scripts.Value, len(values))
	for i, v := range values {
		val, err := scripts.NewValue(v.Type, v.Data)
		if err != nil {
			return nil, err
		}
		jobValues[i] = val
	}
	return jobValues, nil
}

func DTOToFile(_ FileDTO) (*scripts.File, error) {
	return nil, nil
}
