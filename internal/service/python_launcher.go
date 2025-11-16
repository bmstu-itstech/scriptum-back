package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type PythonLauncher struct {
	interpreter string
	directory   string
	maxFileSize int64
	flags       []string

	interpreterToRun string
	installMutex     sync.Mutex
	l                *slog.Logger
}

func NewPythonLauncher(interpreter string, dir string, maxFileSize int64, l *slog.Logger, flags ...string) (*PythonLauncher, error) {
	return &PythonLauncher{
		interpreter:      interpreter,
		directory:        dir,
		maxFileSize:      maxFileSize,
		flags:            flags,
		interpreterToRun: interpreter,
		l:                l,
	}, nil
}

func (p *PythonLauncher) getInterpreter(ctx context.Context, pythonVersion string) (scripts.URL, error) {
	p.l.Debug("Looking for Python version", "version", pythonVersion)
	p.installMutex.Lock()
	defer p.installMutex.Unlock()

	checkCmd := exec.CommandContext(ctx, "pyenv", "versions", "--bare")
	output, err := checkCmd.Output()
	if err != nil {
		return "", fmt.Errorf("can't get pyenv versions: %w", err)
	}
	p.l.Debug("Available pyenv versions", "versions", string(output))

	found := false
	for _, v := range strings.Fields(string(output)) {
		if v == pythonVersion || strings.HasPrefix(v, pythonVersion) {
			found = true
			pythonVersion = v
			break
		}
	}

	if !found {
		p.l.Info("Installing Python version", "version", pythonVersion)
		installCmd := exec.CommandContext(ctx, "pyenv", "install", "-s", pythonVersion)
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			return "", fmt.Errorf("can't install python version: %w", err)
		}
		p.l.Info("Successfully installed Python version", "version", pythonVersion)
	}

	whichCmd := exec.CommandContext(ctx, "pyenv", "which", "python3")
	whichCmd.Env = append(os.Environ(), "PYENV_VERSION="+pythonVersion)

	pathBytes, err := whichCmd.Output()
	if err != nil {
		return "", fmt.Errorf("can't get python interpreter: %w", err)
	}

	interpreterPath := strings.TrimSpace(string(pathBytes))
	p.l.Debug("Found Python interpreter", "path", interpreterPath)
	return scripts.URL(interpreterPath), nil
}

func (p *PythonLauncher) CreateSandbox(ctx context.Context, mainReader scripts.FileData, extraReaders []scripts.FileData, pythonVersion scripts.PythonVersion) (scripts.URL, error) {
	dirName := generateDirname(p.directory)

	if err := os.MkdirAll(dirName, 0755); err != nil {
		return "", fmt.Errorf("can't create sandbox directory: %w", err)
	}

	mainDst := filepath.Join(dirName, mainReader.Name)
	if err := copyFromReader(mainReader.Reader, mainDst); err != nil {
		_ = os.RemoveAll(dirName)
		return "", fmt.Errorf("can't copy main file: %w", err)
	}

	for _, r := range extraReaders {
		dst := filepath.Join(dirName, r.Name)
		if err := copyFromReader(r.Reader, dst); err != nil {
			_ = os.RemoveAll(dirName)
			return "", fmt.Errorf("can't copy extra file: %w", err)
		}
	}

	p.l.Debug("Sandbox created successfully", "dir", dirName, "mainFile", mainDst)
	return mainDst, nil
}

func (p *PythonLauncher) installVenv(ctx context.Context, dirName string) error {
	venv := filepath.Join(dirName, "venv")
	p.l.Debug("Installing virtual environment",
		"dirName", dirName,
		"venv", venv,
		"interpreter", p.interpreterToRun)

	cmd := exec.CommandContext(ctx, p.interpreterToRun, "-m", "venv", venv)
	cmd.Dir = dirName

	var venvStdout, venvStderr bytes.Buffer
	cmd.Stdout = &venvStdout
	cmd.Stderr = &venvStderr

	if err := cmd.Run(); err != nil {
		p.l.Error("Failed to create virtual environment",
			"error", err,
			"stdout", venvStdout.String(),
			"stderr", venvStderr.String())
		_ = os.RemoveAll(dirName)
		return fmt.Errorf("can't create virtual environment: %w", err)
	}

	reqPath := filepath.Join(dirName, "requirements.txt")
	if _, err := os.Stat(reqPath); err == nil {
		p.l.Debug("Requirements.txt found, installing dependencies")
		pipPath := filepath.Join(venv, "bin", "pip")

		installCmd := exec.CommandContext(ctx, pipPath, "install", "-r", "requirements.txt")
		installCmd.Dir = dirName

		var pipStdout, pipStderr bytes.Buffer
		installCmd.Stdout = &pipStdout
		installCmd.Stderr = &pipStderr

		if err := installCmd.Run(); err != nil {
			p.l.Error("Failed to install requirements",
				"error", err,
				"stdout", pipStdout.String(),
				"stderr", pipStderr.String())
			_ = os.RemoveAll(dirName)
			return fmt.Errorf("can't install requirements: %w", err)
		}
		p.l.Debug("Requirements installed successfully",
			"stdout", pipStdout.String())
	} else {
		p.l.Debug("No requirements.txt found", "path", reqPath)
	}

	return nil
}

func (p *PythonLauncher) Run(ctx context.Context, job *scripts.Job) (scripts.Result, error) {
	targetDir := filepath.Dir(job.URL())

	interpreter, err := p.getInterpreter(ctx, job.PythonVersion().String())
	if err != nil {
		_ = os.RemoveAll(targetDir)
		return scripts.Result{}, fmt.Errorf("can't get interpreter: %w", err)
	}

	p.l.Debug("Using Python interpreter", "interpreter", interpreter)
	p.interpreterToRun = interpreter

	err = p.installVenv(ctx, targetDir)
	if err != nil {
		_ = os.RemoveAll(targetDir)
		return scripts.Result{}, err
	}

	args := []string{filepath.Base(job.URL())}
	values := job.Input()
	rawValues := make([]string, len(values))
	for i, v := range values {
		rawValues[i] = v.String()
	}
	args = append(args, rawValues...)

	venvPython := filepath.Join(targetDir, "venv", "bin", "python")

	p.l.Debug("Running Python script",
		"interpreter", venvPython,
		"args", args,
		"dir", targetDir)

	var stdout, stderr bytes.Buffer
	var exitCode int

	cmd := exec.CommandContext(ctx, venvPython, args...)
	cmd.Dir = targetDir
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
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
			_ = os.RemoveAll(targetDir)
			return scripts.Result{}, err
		}
	}

	var res *scripts.Result
	if exitCode == 0 {
		res, err = scripts.NewSuccessResult(outVals)
		if err != nil {
			_ = os.RemoveAll(targetDir)
			return scripts.Result{}, err
		}
		p.l.Debug("Script executed successfully", "jobID", job.ID())
	} else {
		res = scripts.NewFailureResult(exitCode, stderr.String())
		p.l.Warn("Script execution failed",
			"jobID", job.ID(),
			"exitCode", exitCode,
			"stderr", stderr.String())
	}

	return *res, nil
}

func (p *PythonLauncher) DeleteSandbox(ctx context.Context, path scripts.URL) error {
	dir := filepath.Dir(path)
	p.l.Debug("Deleting sandbox", "path", dir)

	err := os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("%w: %s (%w)", scripts.ErrFileNotFound, path, err)
	}

	p.l.Debug("Sandbox deleted successfully", "path", dir)
	return nil
}

func copyFromReader(reader io.Reader, dst string) error {
	out, err := os.Create(dst)
	if err != nil {
		return err
	}

	_, err = io.Copy(out, reader)
	if err != nil {
		_ = os.Remove(dst)
		return err
	}

	if err := out.Close(); err != nil {
		return err
	}

	return nil
}
