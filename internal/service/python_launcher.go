package service

import (
	"bytes"
	"context"
	"os/exec"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type PythonLauncher struct {
	Interpreter string
	Flags       []string
}

type launchResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Err      error
}

func NewPythonLauncher(interpreter string, flags ...string) (*PythonLauncher, error) {
	if interpreter == "" {
		interpreter = "python3"
	}
	return &PythonLauncher{
		Interpreter: interpreter,
		Flags:       flags,
	}, nil
}

func (p *PythonLauncher) Launch(ctx context.Context, job scripts.Job, scriptFields []scripts.Field) (scripts.Result, error) {
	args := []string{job.Command()}
	values := job.In()
	args = append(args, values.Get()...)

	outCh := make(chan launchResult, 1)

	go func() {
		var stdout, stderr bytes.Buffer
		var exitCode int

		cmd := exec.CommandContext(ctx, p.Interpreter, args...)
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

		outCh <- launchResult{
			Stdout:   stdout.String(),
			Stderr:   stderr.String(),
			ExitCode: exitCode,
			Err:      err,
		}
	}()

	launchRes := <-outCh
	if launchRes.Err != nil {
		return scripts.Result{}, scripts.ErrScriptLaunch
	}

	outVals, err := scripts.ParseOutputValues(launchRes.Stdout, scriptFields)
	if err != nil {
		return scripts.Result{}, err
	}

	outVec, err := scripts.NewVector(outVals)
	if err != nil {
		return scripts.Result{}, err
	}
	errMes := scripts.ErrorMessage(launchRes.Stderr)

	result, err := scripts.NewResult(job, scripts.StatusCode(launchRes.ExitCode), *outVec, &errMes, time.Now())
	if err != nil {
		return scripts.Result{}, err
	}

	return *result, nil
}
