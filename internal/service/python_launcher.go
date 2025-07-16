package service

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type PythonLauncher struct{}

func NewPythonLauncher() (*PythonLauncher, error) {
	return &PythonLauncher{}, nil
}

func (p *PythonLauncher) Launch(_ context.Context, job scripts.Job) (scripts.Result, error) {

	return scripts.Result{}, nil
}
