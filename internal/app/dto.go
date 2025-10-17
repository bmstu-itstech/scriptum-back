package app

import (
	"errors"
	"io"
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
	ID            int32
	OwnerID       int64
	MainFileID    int64
	ExtraFileIDs  []int64
	Name          string
	Desc          string
	PythonVersion string
	Input         []FieldDTO
	Output        []FieldDTO
	Visibility    string
	CreatedAt     time.Time
}

type ValueDTO struct {
	Type string
	Data string
}

type ResultDTO struct {
	Output []ValueDTO
	Code   int32
	ErrMsg *string
}

type JobDTO struct {
	JobID        int64
	OwnerID      int64
	ScriptID     int64
	ScriptName   string
	Input        []ValueDTO
	Expected     []FieldDTO
	Url          string
	State        string
	CreatedAt    time.Time
	FinishedAt   *time.Time
	NeedToNotify bool
	JobResult    *ResultDTO
}

type FileDTO struct {
	Name   string
	Reader io.Reader
}

type UserDTO struct {
	UserID  uint32
	Name    string
	Email   string
	IsAdmin bool
}

type ScriptCreateDTO struct {
	OwnerID           int64
	ScriptName        string
	ScriptDescription string
	PythonVersion     string
	MainFileID        int64
	ExtraFileIDs      []int64
	InFields          []FieldDTO
	OutFields         []FieldDTO
}

type ScriptRunDTO struct {
	ScriptID      uint32
	InParams      []ValueDTO
	PythonVersion string
	NeedToNotify  bool
}

type ScriptUpdateDTO struct {
	InFields          []FieldDTO
	OutFields         []FieldDTO
	ScriptName        string
	ScriptDescription string
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
	extra := make([]int64, len(script.ExtraFileIDs()))
	for i, v := range script.ExtraFileIDs() {
		extra[i] = int64(v)
	}
	return ScriptDTO{
		ID:            int32(script.ID()),
		OwnerID:       int64(script.OwnerID()),
		MainFileID:    int64(script.MainFileID()),
		PythonVersion: script.PythonVersion().String(),
		ExtraFileIDs:  extra,
		Name:          script.Name(),
		Desc:          script.Desc(),
		Input:         input,
		Output:        output,
		Visibility:    script.Visibility().String(),
		CreatedAt:     script.CreatedAt(),
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

	extraFiles := make([]scripts.FileID, len(s.ExtraFileIDs))
	for i, id := range s.ExtraFileIDs {
		extraFiles[i] = scripts.FileID(id)
	}

	script, err := scripts.RestoreScript(
		int64(s.ID),
		s.OwnerID,
		s.Name,
		s.Desc,
		s.Visibility,
		s.PythonVersion,
		input,
		output,
		scripts.FileID(s.MainFileID),
		extraFiles,
		s.CreatedAt,
	)
	return script, err
}

func ResultToDTO(r *scripts.Result) (*ResultDTO, error) {
	if r == nil {
		return nil, nil
	}
	output, err := ValuesToDTO(r.Output())
	if err != nil {
		return nil, err
	}
	return &ResultDTO{
		Output: output,
		Code:   int32(r.Code()),
		ErrMsg: r.ErrorMessage(),
	}, nil
}

func JobToDTO(j scripts.Job, name string) (JobDTO, error) {
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
		finishedAt = nil
	}

	res, err := j.Result()
	if err != nil {
		if !errors.Is(err, scripts.ErrJobIsNotFinished) {
			return JobDTO{}, err
		}
	}

	resDto, err := ResultToDTO(res)
	if err != nil {
		return JobDTO{}, err
	}

	return JobDTO{
		JobID:      int64(j.ID()),
		OwnerID:    int64(j.OwnerID()),
		ScriptID:   int64(j.ScriptID()),
		ScriptName: name,
		Input:      input,
		Expected:   expected,
		Url:        j.URL(),
		State:      j.State().String(),
		CreatedAt:  j.CreatedAt(),
		FinishedAt: finishedAt,
		JobResult:  resDto,
	}, nil
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

func UserToDTO(u scripts.User) (UserDTO, error) {
	return UserDTO{
		UserID:  uint32(u.UserID()),
		Name:    string(u.Name()),
		Email:   string(u.Email()),
		IsAdmin: u.IsAdmin(),
	}, nil
}

func DTOToUser(u UserDTO) (*scripts.User, error) {
	return scripts.NewUser(
		scripts.UserID(u.UserID),
		scripts.Name(u.Name),
		scripts.Email(u.Email),
		u.IsAdmin,
	)
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
