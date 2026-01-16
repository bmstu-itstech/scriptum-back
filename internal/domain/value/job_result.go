package value

type JobResult struct {
	code    ExitCode
	output  []Value
	message *string
}

func NewJobResult(code ExitCode, output []Value, message *string) JobResult {
	if output == nil {
		output = make([]Value, 0)
	}
	return JobResult{
		code:    code,
		output:  output,
		message: message,
	}
}

func NewSuccessJobResult(out []Value) JobResult {
	return NewJobResult(0, out, nil)
}

func NewFailureJobResult(code ExitCode, message string) JobResult {
	return NewJobResult(code, nil, &message)
}

func (r JobResult) Code() ExitCode {
	return r.code
}

func (r JobResult) Output() []Value {
	return r.output
}

func (r JobResult) Message() *string {
	return r.message
}
