package scripts

import "time"

type JobID = uint32

type Job struct {
	JobID     JobID
	UserID    UserID
	In        Vector
	Command   string
	startedAt time.Time
}
