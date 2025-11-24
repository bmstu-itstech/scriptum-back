package watermill

const topicRunJob = "run-job"

type payload struct {
	JobID string `json:"job_id"`
}
