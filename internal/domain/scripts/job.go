package scripts

import (
	"fmt"
	"time"
)

type JobID = uint32

type Job struct {
	jobID        JobID
	userID       UserID
	in           Vector
	command      string
	startedAt    time.Time
	outFields    []Field
	userEmail    Email
	needToNotify bool
}

func (j *Job) JobID() JobID {
	return j.jobID
}

func (j *Job) UserID() UserID {
	return j.userID
}

func (j *Job) In() Vector {
	return j.in
}

func (j *Job) Command() string {
	return j.command
}

func (j *Job) StartedAt() time.Time {
	return j.startedAt
}

func (j *Job) OutFields() []Field {
	return j.outFields
}

func (j *Job) UserEmail() Email {
	return j.userEmail
}

func (j *Job) NeedToNotify() bool {
	return j.needToNotify
}

func NewJob(
	jobID JobID,
	userID UserID,
	in Vector,
	command string,
	startedAt time.Time,
	scriptFields []Field,
	userEmail Email,
	needToNotify bool,
) (*Job, error) {
	if in.Len() == 0 {
		return nil, fmt.Errorf("in: expected vector with at least one elemet, got empty vector  %w", ErrJobInvalid)
	}
	return &Job{
		jobID:        jobID,
		userID:       userID,
		in:           in,
		command:      command,
		startedAt:    startedAt,
		outFields:    scriptFields,
		userEmail:    userEmail,
		needToNotify: needToNotify,
	}, nil
}

func NewEmptyJob(
	jobID JobID,
	userID UserID,
	in Vector,
	command string,
	startedAt time.Time,
) (*Job, error) {
	if in.Len() == 0 {
		return nil, fmt.Errorf("in: expected vector with at least one elemet, got empty vector  %w", ErrJobInvalid)
	}
	return &Job{
		jobID:        jobID,
		userID:       userID,
		in:           in,
		command:      command,
		startedAt:    startedAt,
		outFields:    nil,
		userEmail:    "",
		needToNotify: false,
	}, nil
}
