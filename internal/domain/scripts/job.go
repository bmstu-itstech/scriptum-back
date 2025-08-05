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
	scriptID  ScriptID
	input     []Value
	createdAt time.Time
	expected  []Field
	url       URL
}

func NewJobPrototype(ownerID UserID, scriptID ScriptID, input []Value, expected []Field, url URL) (*JobPrototype, error) {
	if ownerID <= 0 {
		return nil, fmt.Errorf("%w: invalid ownerID", ErrInvalidInput)
	}
	if scriptID <= 0 {
		return nil, fmt.Errorf("%w: invalid scriptID", ErrInvalidInput)
	}

	if len(expected) == 0 {
		return nil, fmt.Errorf("%w: invalid Expected: expected at least one output field", ErrInvalidInput)
	}

	if len(url) == 0 {
		return nil, fmt.Errorf("%w: invalid URL: expected not empty URL", ErrInvalidInput)
	}

	if len(url) > FileURLMaxLen {
		return nil, fmt.Errorf(
			"%w: invalid URL: expected len(url) <= %d, got len(url) = %d",
			ErrInvalidInput, FileURLMaxLen, len(url),
		)
	}

	return &JobPrototype{
		ownerID:   ownerID,
		scriptID:  scriptID,
		input:     input,
		expected:  expected,
		url:       url,
		createdAt: time.Now(),
	}, nil
}

func (p *JobPrototype) OwnerID() UserID {
	return p.ownerID
}

func (p *JobPrototype) ScriptID() ScriptID {
	return p.scriptID
}

func (p *JobPrototype) Input() []Value {
	return p.input[:] // Возврат копии
}

func (p *JobPrototype) Expected() []Field {
	return p.expected[:]
}

func (p *JobPrototype) URL() URL {
	return p.url
}

func (p *JobPrototype) CreatedAt() time.Time {
	return p.createdAt
}

func (p *JobPrototype) Build(id JobID) (*Job, error) {
	if id <= 0 {
		return nil, fmt.Errorf("%w: invalid JobID: expected positive id, got %d", ErrInvalidInput, id)
	}

	return &Job{
		JobPrototype: *p,
		id:           id,
		state:        JobPending,
		result:       nil,
		finishedAt:   nil,
	}, nil
}

type Job struct {
	JobPrototype
	id         JobID
	state      JobState
	result     *Result
	finishedAt *time.Time
}

func RestoreJob(
	id int64,
	ownerID int64,
	scriptID int64,
	state string,
	input []Value,
	expected []Field,
	url URL,
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
			scriptID:  ScriptID(scriptID),
			input:     input,
			expected:  expected,
			url:       url,
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
	if j.state != JobRunning {
		return ErrJobIsNotRunning
	}

	j.result = &res
	now := time.Now()
	j.finishedAt = &now
	j.state = JobFinished

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
