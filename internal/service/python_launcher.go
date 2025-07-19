package service

import (
	"bytes"
	"context"
	"encoding/json"
	"os/exec"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/google/uuid"
)

const maxConcurrent = 10

type PythonLauncher struct {
	interpreter string
	flags       []string
	publisher   message.Publisher
	sem         chan struct{}
}

type launchResult struct {
	Stdout   string
	Stderr   string
	ExitCode int
	Err      error
}

type ScriptFinishedEvent struct {
	result    scripts.Result
	userEmail scripts.Email
}

func NewPythonLauncher(interpreter string, publisher message.Publisher, flags ...string) (*PythonLauncher, error) {
	if interpreter == "" {
		interpreter = "python3"
	}
	return &PythonLauncher{
		interpreter: interpreter,
		flags:       flags,
		publisher:   publisher,
		sem:         make(chan struct{}, maxConcurrent),
	}, nil
}

func (p *PythonLauncher) Launch(ctx context.Context, job scripts.Job, scriptFields []scripts.Field, userEmail scripts.Email) (scripts.Result, error) {
	args := []string{job.Command()}
	values := job.In()
	args = append(args, values.Get()...)

	outCh := make(chan launchResult, 1)

	p.sem <- struct{}{}

	go func() {
		defer func() { <-p.sem }()
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

	event := ScriptFinishedEvent{
		result:    *result,
		userEmail: userEmail,
	}

	payload, err := json.Marshal(event)
	if err == nil {
		msg := message.NewMessage(uuid.NewString(), payload)
		_ = p.publisher.Publish("script-finished", msg)

	}

	return *result, nil
}
