package scripts

type StatusCode = int
type ErrorMessage = string

type Result struct {
	job      Job
	code     StatusCode
	out      Vector
	errorMes *ErrorMessage
}

func (r *Result) Job() *Job {
	return &r.job
}

func (r *Result) Code() StatusCode {
	return r.code
}

func (r *Result) Out() *Vector {
	return &r.out
}

func (r *Result) ErrorMessage() *ErrorMessage {
	return r.errorMes
}

func NewResult(job Job, code StatusCode, out Vector, errorMes *ErrorMessage) (*Result, error) {
	return &Result{
		job:      job,
		code:     code,
		out:      out,
		errorMes: errorMes,
	}, nil
}

func NewResultOK(job Job, out Vector) (*Result, error) {
	return &Result{
		job:      job,
		code:     0,
		out:      out,
		errorMes: nil,
	}, nil
}
