package scripts

import "time"

type StatusCode = int
type ErrorMessage = string

type Result struct {
	job      Job
	code     StatusCode
	out      Vector
	errorMes *ErrorMessage
	closedAt time.Time
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

func (r *Result) ClosedAt() time.Time {
	return r.closedAt
}

func NewResult(job Job, code StatusCode, out Vector, errorMes *ErrorMessage, closedAt time.Time) (*Result, error) {
	return &Result{
		job:      job,
		code:     code,
		out:      out,
		errorMes: errorMes,
		closedAt: closedAt,
	}, nil
}

func NewResultOK(job Job, out Vector, closedAt time.Time) (*Result, error) {
	return &Result{
		job:      job,
		code:     0,
		out:      out,
		errorMes: nil,
		closedAt: closedAt,
	}, nil
}
