package api

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	apiv2 "github.com/bmstu-itstech/scriptum-back/gen/go/api/v2"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto"
	"github.com/bmstu-itstech/scriptum-back/internal/app/dto/request"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

func typeFromAPI(ft apiv2.Type) (string, error) {
	switch ft {
	case apiv2.Type_TYPE_INTEGER:
		return value.IntegerValueType.String(), nil
	case apiv2.Type_TYPE_REAL:
		return value.RealValueType.String(), nil
	case apiv2.Type_TYPE_STRING:
		return value.StringValueType.String(), nil
	}
	return "", fmt.Errorf("unknown type: expected one of ['%s', '%s', '%s'], got '%s'",
		value.IntegerValueType.String(), value.RealValueType.String(), value.StringValueType.String(), ft,
	)
}

func typeToAPI(t string) apiv2.Type {
	switch t {
	case value.IntegerValueType.String():
		return apiv2.Type_TYPE_INTEGER
	case value.RealValueType.String():
		return apiv2.Type_TYPE_REAL
	case value.StringValueType.String():
		return apiv2.Type_TYPE_STRING
	default:
		return apiv2.Type_TYPE_STRING // значение по умолчанию, так как мапперы из domain не должны возвращать ошибку
	}
}

func fieldFromAPI(f *apiv2.Field) (dto.Field, error) {
	ft, err := typeFromAPI(f.Type)
	if err != nil {
		return dto.Field{}, err
	}
	return dto.Field{
		Type: ft,
		Name: f.Name,
		Desc: f.Desc,
		Unit: f.Unit,
	}, nil
}

func fieldToAPI(f dto.Field) *apiv2.Field {
	return &apiv2.Field{
		Type: typeToAPI(f.Type),
		Name: f.Name,
		Desc: f.Desc,
		Unit: f.Unit,
	}
}

func fieldsFromAPI(fs []*apiv2.Field) ([]dto.Field, error) {
	res := make([]dto.Field, len(fs))
	for i, f := range fs {
		field, err := fieldFromAPI(f)
		if err != nil {
			return nil, err
		}
		res[i] = field
	}
	return res, nil
}

func fieldsToAPI(fs []dto.Field) []*apiv2.Field {
	res := make([]*apiv2.Field, len(fs))
	for i, f := range fs {
		res[i] = fieldToAPI(f)
	}
	return res
}

func valueFromAPI(v *apiv2.Value) (dto.Value, error) {
	vt, err := typeFromAPI(v.Type)
	if err != nil {
		return dto.Value{}, err
	}
	return dto.Value{
		Type:   vt,
		String: v.Value,
	}, nil
}

func valueToAPI(v dto.Value) *apiv2.Value {
	return &apiv2.Value{
		Type:  typeToAPI(v.Type),
		Value: v.String,
	}
}

func valuesFromAPI(vs []*apiv2.Value) ([]dto.Value, error) {
	res := make([]dto.Value, len(vs))
	for i, v := range vs {
		val, err := valueFromAPI(v)
		if err != nil {
			return nil, err
		}
		res[i] = val
	}
	return res, nil
}

func valuesToAPI(vs []dto.Value) []*apiv2.Value {
	res := make([]*apiv2.Value, len(vs))
	for i, v := range vs {
		res[i] = valueToAPI(v)
	}
	return res
}

func createBoxRequestFromAPI(req *apiv2.CreateBoxRequest, uid int64) (request.CreateBox, error) {
	input, err := fieldsFromAPI(req.Input)
	if err != nil {
		return request.CreateBox{}, err
	}
	output, err := fieldsFromAPI(req.Output)
	if err != nil {
		return request.CreateBox{}, err
	}
	return request.CreateBox{
		UID:       uid,
		ArchiveID: req.ArchiveId,
		Name:      req.Name,
		Desc:      req.Desc,
		Input:     input,
		Output:    output,
	}, nil
}

func visibilityToAPI(s string) apiv2.Visibility {
	switch s {
	case value.VisibilityPrivate.String():
		return apiv2.Visibility_VISIBILITY_PUBLIC
	case value.VisibilityPublic.String():
		return apiv2.Visibility_VISIBILITY_PRIVATE
	default:
		return apiv2.Visibility_VISIBILITY_PRIVATE
	}
}

func boxToAPI(box dto.Box) *apiv2.Box {
	return &apiv2.Box{
		Id:        box.ID,
		OwnerId:   box.OwnerID,
		ArchiveId: box.ArchiveID,
		Name:      box.Name,
		Desc:      box.Desc,
		Vis:       visibilityToAPI(box.Vis),
		In:        fieldsToAPI(box.In),
		Out:       fieldsToAPI(box.Out),
		CreatedAt: timestamppb.New(box.CreatedAt),
	}
}

func boxesToAPI(boxes []dto.Box) []*apiv2.Box {
	res := make([]*apiv2.Box, len(boxes))
	for i, b := range boxes {
		res[i] = boxToAPI(b)
	}
	return res
}

func jobStateToAPI(js string) apiv2.JobState {
	switch js {
	case value.JobPending.String():
		return apiv2.JobState_JOB_STATE_PENDING
	case value.JobRunning.String():
		return apiv2.JobState_JOB_STATE_RUNNING
	case value.JobFinished.String():
		return apiv2.JobState_JOB_STATE_FINISHED
	default:
		return apiv2.JobState_JOB_STATE_PENDING
	}
}

func jobStateFromAPI(js apiv2.JobState) string {
	switch js {
	case apiv2.JobState_JOB_STATE_PENDING:
		return value.JobPending.String()
	case apiv2.JobState_JOB_STATE_RUNNING:
		return value.JobRunning.String()
	case apiv2.JobState_JOB_STATE_FINISHED:
		return value.JobFinished.String()
	default:
		return value.JobPending.String()
	}
}

func jobStateFromAPIOrNil(opt *apiv2.JobState) *string {
	if opt == nil {
		return nil
	}
	js := jobStateFromAPI(*opt)
	return &js
}

func jobResultToAPI(jr dto.JobResult) apiv2.JobResult {
	return apiv2.JobResult{
		Code:    int64(jr.Code),
		Output:  valuesToAPI(jr.Output),
		Message: jr.Message,
	}
}

func jobResultToAPIOrNil(opt *dto.JobResult) *apiv2.JobResult {
	if opt == nil {
		return nil
	}
	jr := jobResultToAPI(*opt)
	return &jr
}

func jobToAPI(job dto.Job) *apiv2.Job {
	return &apiv2.Job{
		Id:         job.ID,
		BoxId:      job.BoxID,
		ArchiveId:  job.ArchiveID,
		OwnerId:    job.OwnerID,
		State:      jobStateToAPI(job.State),
		Input:      valuesToAPI(job.Input),
		Out:        fieldsToAPI(job.Out),
		CreatedAt:  timestamppb.New(job.CreatedAt),
		StartedAt:  timestampbpOrNil(job.StartedAt),
		Result:     jobResultToAPIOrNil(job.Result),
		FinishedAt: timestampbpOrNil(job.FinishedAt),
	}
}

func jobsToAPI(jobs []dto.Job) []*apiv2.Job {
	res := make([]*apiv2.Job, len(jobs))
	for i, j := range jobs {
		res[i] = jobToAPI(j)
	}
	return res
}

func timestampbpOrNil(opt *time.Time) *timestamppb.Timestamp {
	if opt == nil {
		return nil
	}
	return timestamppb.New(*opt)
}

func startJobFromAPI(req *apiv2.StartJobRequest, uid int64) (request.StartJob, error) {
	input, err := valuesFromAPI(req.Values)
	if err != nil {
		return request.StartJob{}, err
	}
	return request.StartJob{
		UID:    uid,
		BoxID:  req.BoxId,
		Values: input,
	}, nil
}
