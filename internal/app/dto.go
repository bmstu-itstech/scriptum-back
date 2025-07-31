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
	Expected     []FieldDTO
	Url          string
	State        string
	CreatedAt    time.Time
	FinishedAt   *time.Time
	NeedToNotify bool
}

func JobToDTO(j scripts.Job) (JobDTO, error) {
	input, err := ValuesToDTO(j.Input())
	if err != nil {
		return JobDTO{}, err
	}

	expected, err := FieldsToDTO(j.Expected())
	if err != nil {
		return JobDTO{}, err
	}

	finishedAt, err := j.FinishedAt()
	if err != nil {
		return JobDTO{}, err
	}
	return JobDTO{
		JobID:      int64(j.ID()),
		OwnerID:    int64(j.OwnerID()),
		ScriptID:   int64(j.ScriptID()),
		Input:      input,
		Expected:   expected,
		Url:        j.URL(),
		State:      j.State().String(),
		CreatedAt:  j.CreatedAt(),
		FinishedAt: finishedAt,
	}, nil

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

func FieldsToDTO(fields []scripts.Field) ([]FieldDTO, error) {
	result := make([]FieldDTO, len(fields))
	for i, field := range fields {
		result[i] = FieldDTO{
			Type: field.ValueType().String(),
			Name: field.Name(),
			Desc: field.Description(),
			Unit: field.Unit(),
		}
	}
	return result, nil
}

func DTOToFields(fields []FieldDTO) ([]scripts.Field, error) {
	jobFields := make([]scripts.Field, len(fields))
	for i, v := range fields {
		valueType, err := scripts.NewValueType(v.Type)
		if err != nil {
			return nil, err
		}
		f, err := scripts.NewField(*valueType, v.Name, v.Desc, v.Unit)
		if err != nil {
			return nil, err
		}
		jobFields[i] = *f
	}
	return jobFields, nil
}

func ScriptToDTO(script scripts.Script) (ScriptDTO, error) {
	input, err := FieldsToDTO(script.Input())
	if err != nil {
		return ScriptDTO{}, nil
	}
	output, err := FieldsToDTO(script.Output())
	if err != nil {
		return ScriptDTO{}, nil
	}
	return ScriptDTO{
		ID:         int32(script.ID()),
		OwnerID:    int64(script.OwnerID()),
		Name:       script.Name(),
		Desc:       script.Desc(),
		Input:      input,
		Output:     output,
		URL:        script.URL(),
		Visibility: script.Visibility().String(),
		CreatedAt:  script.CreatedAt(),
	}, nil
}

func DTOToScript(s ScriptDTO) (*scripts.Script, error) {
	input, err := DTOToFields(s.Input)
	if err != nil {
		return nil, err
	}

	output, err := DTOToFields(s.Output)
	if err != nil {
		return nil, err
	}

	script, err := scripts.RestoreScript(
		int64(s.ID),
		s.OwnerID,
		s.Name,
		s.Desc,
		s.Visibility,
		input,
		output,
		s.URL,
		s.CreatedAt,
	)
	return script, err
}

func DTOToJob(j JobDTO) (*scripts.Job, error) {
	values, err := DTOToValues(j.Input)
	if err != nil {
		return nil, err
	}

	expected, err := DTOToFields(j.Expected)
	if err != nil {
		return nil, err
	}

	job, err := scripts.RestoreJob(
		j.JobID,
		j.OwnerID,
		j.ScriptID,
		j.State,
		values,
		expected,
		j.Url,
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
