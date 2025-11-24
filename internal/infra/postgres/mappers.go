package postgres

import (
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

func boxFieldRowToDomain(row boxFieldRow) (value.Field, error) {
	t, err := value.TypeFromString(row.Type)
	if err != nil {
		return value.Field{}, err
	}
	return value.NewField(t, row.Name, row.Desc, row.Unit)
}

func boxFieldRowsToDomain(rows []boxFieldRow) ([]value.Field, error) {
	res := make([]value.Field, len(rows))
	for i, row := range rows {
		r, err := boxFieldRowToDomain(row)
		if err != nil {
			return nil, err
		}
		res[i] = r
	}
	return res, nil
}

func boxRowToDomain(rBox boxRow, rInput []boxFieldRow, rOutput []boxFieldRow) (*entity.Box, error) {
	in, err := boxFieldRowsToDomain(rInput)
	if err != nil {
		return nil, err
	}
	out, err := boxFieldRowsToDomain(rOutput)
	if err != nil {
		return nil, err
	}
	vis, err := value.VisibilityFromString(rBox.Vis)
	if err != nil {
		return nil, err
	}
	return entity.RestoreBox(
		value.BoxID(rBox.ID),
		value.UserID(rBox.OwnerID),
		value.FileID(rBox.ArchiveID),
		rBox.Name,
		rBox.Desc,
		vis,
		in,
		out,
		rBox.CreatedAt,
	)
}

func boxFieldRowsFromDomain(fields []value.Field, boxID value.BoxID) []boxFieldRow {
	res := make([]boxFieldRow, len(fields))
	for i, field := range fields {
		res[i] = boxFieldRow{
			BoxID: string(boxID),
			Index: i,
			Type:  field.Type().String(),
			Name:  field.Name(),
			Desc:  field.Desc(),
			Unit:  field.Unit(),
		}
	}
	return res
}

func boxRowFromDomain(box *entity.Box) boxRow {
	return boxRow{
		ID:        string(box.ID()),
		OwnerID:   int64(box.OwnerID()),
		ArchiveID: string(box.ArchiveID()),
		Name:      box.Name(),
		Desc:      box.Desc(),
		Vis:       box.Vis().String(),
		CreatedAt: box.CreatedAt(),
	}
}

func jobValueRowToDomain(row jobValueRow) (value.Value, error) {
	t, err := value.TypeFromString(row.Type)
	if err != nil {
		return value.Value{}, err
	}
	return value.NewValue(t, row.Value)
}

func jobValueRowsToDomain(rows []jobValueRow) ([]value.Value, error) {
	res := make([]value.Value, len(rows))
	for i, row := range rows {
		v, err := jobValueRowToDomain(row)
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func jobFieldRowToDomain(row jobFieldRow) (value.Field, error) {
	t, err := value.TypeFromString(row.Type)
	if err != nil {
		return value.Field{}, err
	}
	return value.NewField(t, row.Name, row.Desc, row.Unit)
}

func jobFieldRowsToDomain(rows []jobFieldRow) ([]value.Field, error) {
	res := make([]value.Field, len(rows))
	for i, row := range rows {
		v, err := jobFieldRowToDomain(row)
		if err != nil {
			return nil, err
		}
		res[i] = v
	}
	return res, nil
}

func jobRowToDomain(
	rJob jobRow, rInput []jobValueRow, rOutput []jobValueRow, rOut []jobFieldRow,
) (*entity.Job, error) {
	input, err := jobValueRowsToDomain(rInput)
	if err != nil {
		return nil, err
	}
	output, err := jobValueRowsToDomain(rOutput)
	if err != nil {
		return nil, err
	}
	out, err := jobFieldRowsToDomain(rOut)
	if err != nil {
		return nil, err
	}
	state, err := value.JobStateFromString(rJob.State)
	if err != nil {
		return nil, err
	}
	var result *value.JobResult
	if rJob.ResultCode != nil {
		r := value.NewJobResult(value.ExitCode(*rJob.ResultCode), output, rJob.ResultMsg)
		result = &r
	}
	return entity.RestoreJob(
		value.JobID(rJob.ID),
		value.BoxID(rJob.BoxID),
		value.FileID(rJob.ArchiveID),
		value.UserID(rJob.OwnerID),
		state,
		input,
		out,
		rJob.CreatedAt,
		rJob.StartedAt,
		result,
		rJob.FinishedAt,
	)
}

func jobValueRowsFromDomain(values []value.Value, jobID value.JobID) []jobValueRow {
	res := make([]jobValueRow, len(values))
	for i, v := range values {
		res[i] = jobValueRow{
			JobID: string(jobID),
			Index: i,
			Type:  v.Type().String(),
			Value: v.String(),
		}
	}
	return res
}

func jobFieldRowsFromDomain(fields []value.Field, jobID value.JobID) []jobFieldRow {
	res := make([]jobFieldRow, len(fields))
	for i, f := range fields {
		res[i] = jobFieldRow{
			JobID: string(jobID),
			Index: i,
			Type:  f.Type().String(),
			Name:  f.Name(),
			Desc:  f.Desc(),
			Unit:  f.Unit(),
		}
	}
	return res
}

func jobRowFromDomain(job *entity.Job) jobRow {
	var optCode *int
	var optMsg *string
	var optFinAt *time.Time
	if r := job.Result(); r != nil {
		code := int(r.Code())
		optCode = &code
		optMsg = r.Message()
	}
	return jobRow{
		ID:         string(job.ID()),
		BoxID:      string(job.BoxID()),
		ArchiveID:  string(job.ArchiveID()),
		OwnerID:    int64(job.OwnerID()),
		State:      job.State().String(),
		CreatedAt:  job.CreatedAt(),
		StartedAt:  job.StartedAt(),
		ResultCode: optCode,
		ResultMsg:  optMsg,
		FinishedAt: optFinAt,
	}
}
