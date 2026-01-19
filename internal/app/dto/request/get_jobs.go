package request

type GetJobs struct {
	UID   string
	State *string // optional filter
}
