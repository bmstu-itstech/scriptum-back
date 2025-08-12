package scripts

import (
	"fmt"
)

type StatusCode = int

type Result struct {
	output []Value
	code   StatusCode
	errMsg *string
}

func (r *Result) Code() StatusCode {
	return r.code
}

func (r *Result) Output() []Value {
	return r.output[:]
}

func (r *Result) ErrorMessage() *string {
	return r.errMsg
}

func RestoreResult(output []Value, code StatusCode, errMsg *string) *Result {
	var msg *string
	var out []Value
	if errMsg != nil {
		msg = errMsg
	}
	if output != nil {
		out = output[:]
	}
	return &Result{
		output: out,
		code:   code,
		errMsg: msg,
	}
}

func NewSuccessResult(output []Value) (*Result, error) {
	if len(output) == 0 {
		// Скрипт не может не иметь выходных значений, поэтому ошибка программиста.
		return nil, fmt.Errorf("empty output")
	}

	return &Result{
		output: output[:], // Копирование
		code:   0,
		errMsg: nil,
	}, nil
}

func NewFailureResult(code StatusCode, errMsg string) *Result {
	return &Result{
		output: nil,
		code:   code,
		errMsg: &errMsg,
	}
}
