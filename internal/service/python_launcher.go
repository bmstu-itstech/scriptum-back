package service

import (
	"bytes"
	"context"
	"os/exec"
	"time"

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

func (p *PythonLauncher) Launch(ctx context.Context, job scripts.Job) (scripts.Result, error) {
	args := []string{job.Command()}
	values := job.In()
	args = append(args, values.Get()...)

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
		return scripts.Result{}, scripts.ErrScriptLaunch
	}

	outVals, err := scripts.ParseOutputValues(stdout.String(), job.ScriptFields())
	if err != nil {
		return scripts.Result{}, err
	}

	outVec, err := scripts.NewVector(outVals)
	if err != nil {
		return scripts.Result{}, err
	}
	errMes := scripts.ErrorMessage(stderr.String())

	result, err := scripts.NewResult(job, scripts.StatusCode(exitCode), *outVec, &errMes, time.Now())
	if err != nil {
		return scripts.Result{}, err
	}

	return *result, nil
}
