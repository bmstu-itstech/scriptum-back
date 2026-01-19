package dto

import (
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/entity"
)

type Job struct {
	ID          string
	BlueprintID string
	ArchiveID   string
	OwnerID     string
	State       string
	Input       []Value
	Out         []Field
	CreatedAt   time.Time
	StartedAt   *time.Time
	Result      *JobResult
	FinishedAt  *time.Time
}

func JobToDTO(job *entity.Job) Job {
	var optRes *JobResult
	if res := job.Result(); res != nil {
		t := JobResultToDTO(*res)
		optRes = &t
	}
	return Job{
		ID:          string(job.ID()),
		BlueprintID: string(job.BlueprintID()),
		ArchiveID:   string(job.ArchiveID()),
		OwnerID:     string(job.OwnerID()),
		State:       job.State().String(),
		Input:       valuesToDTOs(job.Input()),
		Out:         fieldsToDTOs(job.Out()),
		CreatedAt:   job.CreatedAt(),
		StartedAt:   job.StartedAt(),
		Result:      optRes,
		FinishedAt:  job.FinishedAt(),
	}
}

func JobsToDTOs(jobs []*entity.Job) []Job {
	res := make([]Job, len(jobs))
	for i, job := range jobs {
		res[i] = JobToDTO(job)
	}
	return res
}
