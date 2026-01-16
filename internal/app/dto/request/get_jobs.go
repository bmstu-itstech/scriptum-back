package request

type GetJobs struct {
	UID   int64
	State *string // optional filter
}
