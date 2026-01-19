package postgres

import (
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

func blueprintFieldRowToDomain(row blueprintFieldRow) (value.Field, error) {
	t, err := value.TypeFromString(row.Type)
	if err != nil {
		return value.Field{}, err
	}
	return value.NewField(t, row.Name, row.Desc, row.Unit)
}

func blueprintFieldRowsToDomain(rows []blueprintFieldRow) ([]value.Field, error) {
	res := make([]value.Field, len(rows))
	for i, row := range rows {
		r, err := blueprintFieldRowToDomain(row)
		if err != nil {
			return nil, err
		}
		res[i] = r
	}
	return res, nil
}

func blueprintRowToDomain(rB blueprintRow, rInput []blueprintFieldRow, rOutput []blueprintFieldRow) (*entity.Blueprint, error) {
	in, err := blueprintFieldRowsToDomain(rInput)
	if err != nil {
		return nil, err
	}
	out, err := blueprintFieldRowsToDomain(rOutput)
	if err != nil {
		return nil, err
	}
	vis, err := value.VisibilityFromString(rB.Vis)
	if err != nil {
		return nil, err
	}
	return entity.RestoreBlueprint(
		value.BlueprintID(rB.ID),
		value.UserID(rB.OwnerID),
		value.FileID(rB.ArchiveID),
		rB.Name,
		rB.Desc,
		vis,
		in,
		out,
		rB.CreatedAt,
	)
}

func blueprintFieldRowsFromDomain(fields []value.Field, blueprintID value.BlueprintID) []blueprintFieldRow {
	res := make([]blueprintFieldRow, len(fields))
	for i, field := range fields {
		res[i] = blueprintFieldRow{
			BlueprintID: string(blueprintID),
			Index:       i,
			Type:        field.Type().String(),
			Name:        field.Name(),
			Desc:        field.Desc(),
			Unit:        field.Unit(),
		}
	}
	return res
}

func blueprintRowFromDomain(b *entity.Blueprint) blueprintRow {
	return blueprintRow{
		ID:        string(b.ID()),
		OwnerID:   string(b.OwnerID()),
		ArchiveID: string(b.ArchiveID()),
		Name:      b.Name(),
		Desc:      b.Desc(),
		Vis:       b.Vis().String(),
		CreatedAt: b.CreatedAt(),
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
		value.BlueprintID(rJob.BlueprintID),
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

func userRowToDomain(row userRow) (*entity.User, error) {
	role, err := value.RoleFromString(row.Role)
	if err != nil {
		return nil, err
	}
	return entity.RestoreUser(
		value.UserID(row.ID),
		value.MustEmailFromString(row.Email),
		[]byte(row.Passhash),
		row.Name,
		role,
		row.CreatedAt,
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
		ID:          string(job.ID()),
		BlueprintID: string(job.BlueprintID()),
		ArchiveID:   string(job.ArchiveID()),
		OwnerID:     string(job.OwnerID()),
		State:       job.State().String(),
		CreatedAt:   job.CreatedAt(),
		StartedAt:   job.StartedAt(),
		ResultCode:  optCode,
		ResultMsg:   optMsg,
		FinishedAt:  optFinAt,
	}
}
