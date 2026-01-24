package request

type GetJobs struct {
	ActorID string
	State   *string // optional filter
}
