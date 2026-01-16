package value

type Result struct {
	code   ExitCode
	output string
}

func NewResult(code ExitCode) Result {
	return Result{
		code: code,
	}
}

func (r Result) WithOutput(o string) Result {
	r.output = o
	return r
}

func (r Result) Code() ExitCode {
	return r.code
}

func (r Result) Output() string {
	return r.output
}
