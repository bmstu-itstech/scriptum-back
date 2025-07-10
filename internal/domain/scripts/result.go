package scripts

type StatusCode = int
type ErrorMessage = string

type Result struct {
	JobID    Job
	Сode     StatusCode
	out      Vector
	errorMes ErrorMessage
}

func (r Result) Out() Vector {
	return r.out
}

func (r Result) ErrorMessage() ErrorMessage {
	return r.errorMes
}
