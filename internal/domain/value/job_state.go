package value

type JobState struct {
	s string
}

var (
	JobPending  = JobState{"pending"}
	JobRunning  = JobState{"running"}
	JobFinished = JobState{"finished"}
)

func (j JobState) String() string {
	return j.s
}

func (j JobState) IsZero() bool {
	return j.s == ""
}
