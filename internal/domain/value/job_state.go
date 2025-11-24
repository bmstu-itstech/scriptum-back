package value

import (
	"fmt"

	"github.com/bmstu-itstech/scriptum-back/internal/domain"
)

type JobState struct {
	s string
}

var (
	JobPending  = JobState{"pending"}
	JobRunning  = JobState{"running"}
	JobFinished = JobState{"finished"}
)

func JobStateFromString(s string) (JobState, error) {
	switch s {
	case "pending":
		return JobPending, nil
	case "running":
		return JobRunning, nil
	case "finished":
		return JobFinished, nil
	}
	return JobPending, domain.NewInvalidInputError(
		"job-state-invalid",
		fmt.Sprintf("invalid job state: expected one of ['pending', 'running', 'finished'], got '%s'", s),
	)
}

func (j JobState) String() string {
	return j.s
}

func (j JobState) IsZero() bool {
	return j.s == ""
}
