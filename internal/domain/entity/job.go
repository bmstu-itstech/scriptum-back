package entity

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

var ErrInvalidJobStateChange = errors.New("invalid job state change")
var ErrJobResultParseFailed = errors.New("job result parse failed")

type Job struct {
	id          value.JobID
	blueprintID value.BlueprintID
	archiveID   value.FileID
	ownerID     value.UserID
	state       value.JobState
	input       []value.Value
	out         []value.Field
	createdAt   time.Time

	startedAt  *time.Time
	result     *value.JobResult
	finishedAt *time.Time
}

func (j *Job) Run() error {
	if j.state != value.JobPending {
		return fmt.Errorf(
			"%w: expected JobPending -> JobRunning, got %s -> JobRunning", ErrInvalidJobStateChange, j.state.String(),
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
			"%w: expected JobRunning -> JobFinished, got %s -> JobFinished", ErrInvalidJobStateChange, j.state.String(),
		)
	}
	if res.Code().IsSuccess() {
		out, err := j.parseOutput(res.Output())
		if err != nil {
			return err
		}
		jRes := value.NewSuccessJobResult(out)
		j.result = &jRes
	} else {
		jRes := value.NewFailureJobResult(res.Code(), res.Output())
		j.result = &jRes
	}
	j.state = value.JobFinished
	now := time.Now()
	j.finishedAt = &now
	return nil
}

func (j *Job) parseOutput(output string) ([]value.Value, error) {
	lines := strings.Split(output, "\n")
	lines = lines[:len(lines)-1]
	if len(lines) != len(j.out) {
		return nil, fmt.Errorf("%w: expected %d lines, got %d", ErrJobResultParseFailed, len(j.out), len(lines))
	}
	res := make([]value.Value, len(j.out))
	for i, line := range lines {
		field := j.out[i]
		v, err := value.NewValue(field.Type(), line)
		if err != nil {
			return nil, fmt.Errorf("%w: line=%d: %w", ErrJobResultParseFailed, i+1, err)
		}
		res[i] = v
	}
	return res, nil
}

func (j *Job) ID() value.JobID {
	return j.id
}

func (j *Job) BlueprintID() value.BlueprintID {
	return j.blueprintID
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

func (j *Job) Input() []value.Value {
	return j.input
}

func (j *Job) Out() []value.Field {
	return j.out
}

func (j *Job) CreatedAt() time.Time {
	return j.createdAt
}

func (j *Job) StartedAt() *time.Time {
	return j.startedAt
}

func (j *Job) Result() *value.JobResult {
	return j.result
}

func (j *Job) FinishedAt() *time.Time {
	return j.finishedAt
}

func RestoreJob(
	id value.JobID,
	blueprintID value.BlueprintID,
	archiveID value.FileID,
	ownerID value.UserID,
	state value.JobState,
	input []value.Value,
	out []value.Field,
	createdAt time.Time,
	startedAt *time.Time,
	result *value.JobResult,
	finishedAt *time.Time,
) (*Job, error) {
	if id == "" {
		return nil, errors.New("empty id")
	}

	if blueprintID == "" {
		return nil, errors.New("empty blueprintID")
	}

	if archiveID == "" {
		return nil, errors.New("empty archiveID")
	}

	if ownerID == "" {
		return nil, errors.New("empty ownerID")
	}

	if state.IsZero() {
		return nil, errors.New("empty state")
	}

	if input == nil {
		input = make([]value.Value, 0)
	}

	if out == nil {
		out = make([]value.Field, 0)
	}

	return &Job{
		id:          id,
		blueprintID: blueprintID,
		archiveID:   archiveID,
		ownerID:     ownerID,
		state:       state,
		input:       input,
		out:         out,
		createdAt:   createdAt,
		startedAt:   startedAt,
		result:      result,
		finishedAt:  finishedAt,
	}, nil
}
