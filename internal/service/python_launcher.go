package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type PythonLauncher struct {
	interpreter string
	directory   string
	maxFileSize int64
	flags       []string
}

func NewPythonLauncher(interpreter string, dir string, maxFileSize int64, flags ...string) (*PythonLauncher, error) {
	return &PythonLauncher{
		interpreter: interpreter,
		directory:   dir,
		maxFileSize: maxFileSize,
		flags:       flags,
	}, nil
}

func (p *PythonLauncher) CreateSandbox(ctx context.Context, mainReader scripts.FileData, extraReaders []scripts.FileData) (scripts.URL, error) {
	dirName := generateDirname(p.directory)

	if err := os.MkdirAll(dirName, 0755); err != nil {
		return "", err
	}
	mainDst := filepath.Join(dirName, mainReader.Name)
	if err := copyFromReader(mainReader.Reader, mainDst, p.maxFileSize); err != nil {
		_ = os.RemoveAll(dirName)
		return "", err
	}

	for _, r := range extraReaders {
		dst := filepath.Join(dirName, r.Name)

		if err := copyFromReader(r.Reader, dst, p.maxFileSize); err != nil {
			_ = os.RemoveAll(dirName)
			return "", err
		}
	}

	return mainDst, nil
}

func (p *PythonLauncher) Run(ctx context.Context, job *scripts.Job) (scripts.Result, error) {
	args := []string{job.URL()}
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

	out := strings.Fields(stdout.String())

	expected := job.Expected()

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

	return *res, nil
}

func (p *PythonLauncher) DeleteSandbox(ctx context.Context, path scripts.URL) error {
	dir := filepath.Dir(path)

	err := os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("%w: %s (%w)", scripts.ErrFileNotFound, path, err)
	}
	return nil
}

func copyFromReader(reader io.Reader, dst string, maxSize int64) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	limited := io.LimitReader(reader, maxSize+1)
	written, err := io.Copy(out, limited)
	if err != nil {
		_ = os.Remove(dst)
		return err
	}

	if written > maxSize {
		_ = os.Remove(dst)
		return fmt.Errorf("%w: file size exceeds limit of %d bytes", scripts.ErrFileUpload, maxSize)
	}

	return nil
}
