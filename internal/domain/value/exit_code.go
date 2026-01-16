package value

type ExitCode int

func (e ExitCode) IsSuccess() bool {
	return e == 0
}
