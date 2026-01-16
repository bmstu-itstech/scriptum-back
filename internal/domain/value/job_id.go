package value

const JobIDLength = 8

type JobID string

func NewJobID() JobID {
	return JobID(NewShortUUID(JobIDLength))
}
