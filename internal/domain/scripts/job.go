package scripts

import (
	"errors"
	"fmt"
	"time"
)

type JobID int64 // JobID > 0

type JobState struct {
	s string
}

func (j JobState) String() string {
	return j.s
}

var JobPending = JobState{"pending"}
var JobRunning = JobState{"running"}
var JobFinished = JobState{"finished"}

func NewJobStateFromString(s string) (JobState, error) {
	switch s {
	case "pending":
		return JobPending, nil
	case "running":
		return JobRunning, nil
	case "finished":
		return JobFinished, nil
	}
	return JobState{}, fmt.Errorf(
		"%w: invalid JobState: expected one of ['pending', 'running', 'finished'], got %s",
		ErrInvalidInput, s,
	)
}

type JobPrototype struct {
	ownerID   UserID
	input     []Value
	createdAt time.Time
}

func (p *JobPrototype) OwnerID() UserID {
	return p.ownerID
}

func (p *JobPrototype) Input() []Value {
	return p.input[:] // Возврат копии
}

func (p *JobPrototype) CreatedAt() time.Time {
	return p.createdAt
}

func (p *JobPrototype) Build(ownerID UserID, id JobID) (*Job, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid JobID: expected positive id, got %d", ErrInvalidInput, id)
	}

	if ownerID == 0 {
		return nil, fmt.Errorf("empty owner ID ")
	}

	return &Job{
		JobPrototype: *p,
		id:           id,
		ownerID:      ownerID,
		state:        JobPending,
		result:       nil,
		finishedAt:   nil,
	}, nil
}

type Job struct {
	JobPrototype
	id         JobID
	ownerID    UserID
	state      JobState
	result     *Result
	finishedAt *time.Time
}

func RestoreJob(
	id int64,
	ownerID int64,
	state string,
	input []Value,
	result *Result,
	createdAt time.Time,
	finishedAt *time.Time,
) (*Job, error) {
	if id == 0 {
		return nil, fmt.Errorf("job.id is empty")
	}

	if state == "" {
		return nil, fmt.Errorf("job.state is empty")
	}

	jState, err := NewJobStateFromString(state)
	if err != nil {
		return nil, fmt.Errorf("invalid job.state %s", state)
	}

	return &Job{
		JobPrototype: JobPrototype{
			ownerID:   UserID(ownerID),
			input:     input,
			createdAt: createdAt,
		},
		id:         JobID(id),
		state:      jState,
		result:     result,
		finishedAt: finishedAt,
	}, nil
}

var ErrJobIsNotPending = errors.New("job is not pending yet")

func (j *Job) Run() error {
	if j.state != JobPending {
		return ErrJobIsNotPending
	}

	j.state = JobRunning

	return nil
}

var ErrJobIsNotRunning = errors.New("job is not running")

func (j *Job) Finish(res Result) error {
	if j.state == JobFinished {
		return ErrJobIsNotRunning
	}

	j.result = &res
	now := time.Now()
	j.finishedAt = &now

	return nil
}

func (j *Job) ID() JobID {
	return j.id
}

func (j *Job) State() JobState {
	return j.state
}

var ErrJobIsNotFinished = errors.New("job is not finished")

func (j *Job) Result() (*Result, error) {
	if j.state != JobFinished {
		return nil, ErrJobIsNotFinished
	}
	return j.result, nil
}

func (j *Job) FinishedAt() (*time.Time, error) {
	if j.state != JobFinished {
		return nil, ErrJobIsNotFinished
	}
	return j.finishedAt, nil
}

func (j *Job) Duration() (time.Duration, error) {
	if j.state != JobFinished {
		return 0, ErrJobIsNotFinished
	}
	return j.finishedAt.Sub(j.JobPrototype.createdAt), nil
}
