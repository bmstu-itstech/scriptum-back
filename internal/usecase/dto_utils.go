package usecase

import (
	"strconv"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type FieldDTO struct {
	Type        string
	Name        string
	Description string
	Unit        string
}

type ScriptDTO struct {
	ID         uint32
	Fields     []FieldDTO
	Path       string
	Owner      uint32
	Visibility string
	CreatedAt  time.Time
}

type ValueDTO struct {
	Type string
	Data string
}

type JobDTO struct {
	JobID     uint32
	UserID    uint32
	In        []ValueDTO
	Command   string
	StartedAt time.Time
}

type ResultDTO struct {
	Job      JobDTO
	Code     int
	Out      []ValueDTO
	ErrorMes *string
}

func FieldsToDTO(fields []scripts.Field) []FieldDTO {
	dto := make([]FieldDTO, 0, len(fields))
	for _, field := range fields {
		sType := field.FieldType()
		dto = append(dto, FieldDTO{
			Type:        sType.String(),
			Name:        field.Name(),
			Description: field.Description(),
			Unit:        field.Unit(),
		})
	}
	return dto
}

func DTOToFields(dto []FieldDTO) ([]scripts.Field, error) {
	fields := make([]scripts.Field, 0, len(dto))
	for _, f := range dto {
		type_, err := scripts.NewType(f.Type)
		if err != nil {
			return nil, err
		}
		field, err := scripts.NewField(*type_, f.Name, f.Description, f.Unit)
		if err != nil {
			return nil, err
		}
		fields = append(fields, *field)
	}
	return fields, nil
}

func ScriptToDTO(script scripts.Script) ScriptDTO {
	fds := FieldsToDTO(script.Fields())
	return ScriptDTO{
		ID:         uint32(script.ID()),
		Fields:     fds,
		Path:       script.Path(),
		Owner:      uint32(script.Owner()),
		Visibility: string(script.Visibility()),
		CreatedAt:  script.CreatedAt(),
	}
}

func DTOToScript(dto ScriptDTO) (scripts.Script, error) {
	fields, err := DTOToFields(dto.Fields)
	if err != nil {
		return scripts.Script{}, err
	}
	res, err := scripts.NewScript(
		fields,
		scripts.Path(dto.Path),
		scripts.UserID(dto.Owner),
		scripts.Visibility(dto.Visibility),
	)
	return *res, err
}

func VectorToDTO(v scripts.Vector) []ValueDTO {
	dto := make([]ValueDTO, 0, v.Len())
	for _, val := range v.Values() {
		valType := val.VariableType()
		dto = append(dto, ValueDTO{
			Type: valType.String(),
			Data: val.String(),
		})
	}
	return dto
}

func DTOToVector(dto []ValueDTO) (scripts.Vector, error) {
	valuesVec := make([]scripts.Value, 0, len(dto))
	for _, v := range dto {
		switch v.Type {
		case "integer":
			data, err := strconv.ParseInt(v.Data, 10, 64)
			if err != nil {
				return scripts.Vector{}, err
			}

			val, err := scripts.NewInteger(data)
			if err != nil {
				return scripts.Vector{}, err
			}

			valuesVec = append(valuesVec, val)

		case "real":
			data, err := strconv.ParseFloat(v.Data, 64)
			if err != nil {
				return scripts.Vector{}, err
			}

			val, err := scripts.NewReal(data)
			if err != nil {
				return scripts.Vector{}, err
			}

			valuesVec = append(valuesVec, val)

		case "complex":
			data, err := strconv.ParseComplex(v.Data, 64)
			if err != nil {
				return scripts.Vector{}, err
			}

			val, err := scripts.NewComplex(complex64(data))
			if err != nil {
				return scripts.Vector{}, err
			}

			valuesVec = append(valuesVec, val)

		default:
			return scripts.Vector{}, scripts.ErrInvalidValueType
		}
	}
	vector, _ := scripts.NewVector(valuesVec)
	return *vector, nil
}

func JobToDTO(job scripts.Job) JobDTO {
	dto := JobDTO{
		JobID:     uint32(job.JobID()),
		UserID:    uint32(job.UserID()),
		In:        VectorToDTO(job.In()),
		Command:   job.Command(),
		StartedAt: job.StartedAt(),
	}
	return dto
}

func DTOToJob(dto JobDTO) (scripts.Job, error) {
	in, err := DTOToVector(dto.In)
	if err != nil {
		return scripts.Job{}, err
	}

	job, err := scripts.NewJob(
		scripts.JobID(dto.JobID),
		scripts.UserID(dto.UserID),
		in,
		dto.Command,
		dto.StartedAt,
	)
	if err != nil {
		return scripts.Job{}, err
	}

	return *job, nil
}

func ResultToDTO(result scripts.Result) ResultDTO {
	return ResultDTO{
		Job:      JobToDTO(*result.Job()),
		Code:     result.Code(),
		Out:      VectorToDTO(*result.Out()),
		ErrorMes: result.ErrorMessage(),
	}
}

func DTOToResult(dto ResultDTO) (scripts.Result, error) {
	job, err := DTOToJob(dto.Job)
	if err != nil {
		return scripts.Result{}, err
	}

	out, err := DTOToVector(dto.Out)
	if err != nil {
		return scripts.Result{}, err
	}

	result, err := scripts.NewResult(
		job,
		dto.Code,
		out,
		dto.ErrorMes,
	)
	if err != nil {
		return scripts.Result{}, err
	}

	return *result, nil
}
