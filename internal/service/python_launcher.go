package service

import (
	"bytes"
	"context"
	"os/exec"
	"strings"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type PythonLauncher struct {
	interpreter string
	flags       []string
	publisher   message.Publisher
}

func NewPythonLauncher(interpreter string, publisher message.Publisher, flags ...string) (*PythonLauncher, error) {
	if interpreter == "" {
		interpreter = "python3"
	}
	return &PythonLauncher{
		interpreter: interpreter,
		flags:       flags,
		publisher:   publisher,
	}, nil
}

func (p *PythonLauncher) Run(ctx context.Context, job *scripts.Job, path string, expected []scripts.Field) (scripts.Result, error) {
	args := []string{path}
	values := job.Input()
	rawValues := make([]string, len(values))
	for i, v := range values {
		rawValues[i] = v.String()
	}
	args = append(args, rawValues...)

	var stdout, stderr bytes.Buffer
	var exitCode int

	cmd := exec.CommandContext(ctx, p.interpreter, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if exitErr, ok := err.(*exec.ExitError); ok {
		exitCode = exitErr.ExitCode()
	} else if err == nil {
		exitCode = 0
	} else {
		exitCode = -1
	}

	if err != nil {
		return scripts.Result{}, err
	}

	out := strings.Fields(stdout.String())

	outVals := make([]scripts.Value, len(expected))
	for i, token := range out {
		outVals[i], err = scripts.NewValue(expected[i].ValueType().String(), token)
		if err != nil {
			return scripts.Result{}, err
		}
	}

	var res *scripts.Result
	if exitCode == 0 {
		res, err = scripts.NewSuccessResult(outVals)
		if err != nil {
			return scripts.Result{}, err
		}
	} else {
		res = scripts.NewFailureResult(exitCode, stderr.String())
	}
	if err != nil {
		return scripts.Result{}, err
	}

	return *res, nil
}
