package entities

import (
	"errors"
	"fmt"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

var InvalidJobStateChange = errors.New("invalid job state change")

type Job struct {
	id         value.JobID
	boxID      value.BoxID
	archiveID  value.FileID
	ownerID    value.UserID
	state      value.JobState
	input      value.Input
	result     *value.Result
	createdAt  time.Time
	startedAt  *time.Time
	finishedAt *time.Time
}

func (j *Job) Run() error {
	if j.state != value.JobPending {
		return fmt.Errorf(
			"%w: expected JobPending -> JobRunning, got %s -> JobRunning", InvalidJobStateChange, j.state.String(),
		)
	}
	j.state = value.JobRunning
	now := time.Now()
	j.startedAt = &now
	return nil
}

func (j *Job) Finish(res value.Result) error {
	if j.state != value.JobRunning {
		return fmt.Errorf(
			"%w: expected JobRunning -> JobFinished, got %s -> JobFinished", InvalidJobStateChange, j.state.String(),
		)
	}
	j.state = value.JobFinished
	now := time.Now()
	j.finishedAt = &now
	j.result = &res
	return nil
}

func (j *Job) ID() value.JobID {
	return j.id
}

func (j *Job) BoxID() value.BoxID {
	return j.boxID
}

func (j *Job) ArchiveID() value.FileID {
	return j.archiveID
}

func (j *Job) OwnerID() value.UserID {
	return j.ownerID
}

func (j *Job) State() value.JobState {
	return j.state
}

func (j *Job) Input() value.Input {
	return j.input
}

func (j *Job) CreatedAt() time.Time {
	return j.createdAt
}

func (j *Job) StartedAt() *time.Time {
	return j.startedAt
}

func (j *Job) FinishedAt() *time.Time {
	return j.finishedAt
}

func (j *Job) Result() *value.Result {
	return j.result
}
