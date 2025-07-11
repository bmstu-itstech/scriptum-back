package scripts

import "time"

type JobID = uint32

type Job struct {
	jobID     JobID
	userID    UserID
	in        Vector
	command   string
	startedAt time.Time
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

func NewJob(jobID JobID, userID UserID, in Vector, command string) (*Job, error) {
	if jobID == 0 {
		return nil, ErrInvalidJobID
	}
	if userID == 0 {
		return nil, ErrInvalidUserID
	}
	if len(in.Values()) == 0 {
		return nil, ErrEmptyVector
	}
	if command == "" {
		return nil, ErrInvalidCommand
	}

	return &Job{
		jobID:     jobID,
		userID:    userID,
		in:        in,
		command:   command,
		startedAt: time.Now(),
	}, nil
}
