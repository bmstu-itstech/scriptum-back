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

type ResultDTO struct {
	Code   int
	Output []ValueDTO
	ErrMsg *string
}

type JobDTO struct {
	JobID      int64
	OwnerID    int64
	Input      []ValueDTO
	State      string
	Result     *ResultDTO
	CreatedAt  time.Time
	FinishedAt *time.Time
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

func DTOToJob(_ JobDTO) (*scripts.Job, error) {
	return nil, nil
}

func DTOToValues(_ []ValueDTO) ([]scripts.Value, error) {
	return nil, nil
}

func DTOToFile(_ FileDTO) (*scripts.File, error) {
	return nil, nil
}
