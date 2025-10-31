package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
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

	interpreterToRun string
}

func NewPythonLauncher(interpreter string, dir string, maxFileSize int64, flags ...string) (*PythonLauncher, error) {
	return &PythonLauncher{
		interpreter: interpreter,
		directory:   dir,
		maxFileSize: maxFileSize,
		flags:       flags,

		interpreterToRun: interpreter,
	}, nil
}

func (p *PythonLauncher) getInterpreter(ctx context.Context, pythonVersion string) (scripts.URL, error) {
	checkCmd := exec.CommandContext(ctx, "pyenv", "versions", "--bare")
	output, err := checkCmd.Output()
	if err != nil {
		return "", fmt.Errorf("can't get pyenv versions: %w", err)
	}

	found := false
	for _, v := range strings.Fields(string(output)) {
		if v == pythonVersion || strings.HasPrefix(v, pythonVersion) {
			found = true
			break
		}
	}

	if !found {
		installCmd := exec.CommandContext(ctx, "pyenv", "install", "-s", pythonVersion)
		installCmd.Stdout = os.Stdout
		installCmd.Stderr = os.Stderr
		if err := installCmd.Run(); err != nil {
			return "", fmt.Errorf("can't install python version: %w", err)
		}
	}

	whichCmd := exec.CommandContext(ctx, "pyenv", "which", "python3")
	whichCmd.Env = append(os.Environ(), "PYENV_VERSION="+pythonVersion)

	pathBytes, err := whichCmd.Output()
	if err != nil {
		return "", fmt.Errorf("can't get python interpreter: %w", err)
	}

	return scripts.URL(strings.TrimSpace(string(pathBytes))), nil
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

	// venv := filepath.Join(dirName, "venv")
	// log.Println("CreateSandbox dirName", dirName)
	// log.Println("CreateSandbox venv", venv)

	interpreter, err := p.getInterpreter(ctx, pythonVersion.String())
	if err != nil {
		_ = os.RemoveAll(dirName)
		return "", fmt.Errorf("can't get interpreter: %w", err)
	}
	p.interpreterToRun = interpreter
	log.Println("CreateSandbox interpreter", interpreter)

	// cmd := exec.CommandContext(ctx, interpreter, "-m", "venv", venv)
	// cmd.Dir = dirName
	// log.Println("command: ", interpreter, "-m", "venv", venv)
	// if err := cmd.Run(); err != nil {
	// 	_ = os.RemoveAll(dirName)
	// 	return "", fmt.Errorf("can't create virtual environment: %w", err)
	// }

	// reqPath := filepath.Join(dirName, "requirements.txt")

	// if _, err := os.Stat(reqPath); err == nil {
	// 	log.Println("requirements.txt found")
	// 	pipPath := filepath.Join(venv, "bin", "pip")
	// 	installCmd := exec.CommandContext(ctx, pipPath, "install", "-r", "requirements.txt")
	// 	installCmd.Dir = dirName
	// 	if err := installCmd.Run(); err != nil {
	// 		log.Println("requirements.txt found, but can't install")
	// 		_ = os.RemoveAll(dirName)
	// 		return "", fmt.Errorf("can't install requirements: %w", err)
	// 	}
	// 	log.Println("requirements installed")
	// }

	return mainDst, nil
}

func (p *PythonLauncher) installVenv(ctx context.Context, dirName string) error {
	venv := filepath.Join(dirName, "venv")
	log.Println("installVenv dirName", dirName)
	log.Println("installVenv", venv)

	interpreter := p.interpreterToRun
	log.Println("installVenv interpreter", interpreter)

	cmd := exec.CommandContext(ctx, interpreter, "-m", "venv", venv)
	cmd.Dir = dirName
	log.Println(interpreter, "-m", "venv", venv)
	if err := cmd.Run(); err != nil {
		_ = os.RemoveAll(dirName)
		return fmt.Errorf("can't create virtual environment: %w", err)
	}

	reqPath := filepath.Join(dirName, "requirements.txt")

	if _, err := os.Stat(reqPath); err == nil {
		log.Println("requirements.txt found")
		pipPath := filepath.Join(venv, "bin", "pip")
		installCmd := exec.CommandContext(ctx, pipPath, "install", "-r", "requirements.txt")
		installCmd.Dir = dirName
		if err := installCmd.Run(); err != nil {
			log.Println("requirements.txt found, but can't install", err)
			_ = os.RemoveAll(dirName)
			return fmt.Errorf("can't install requirements: %w", err)
		}
		log.Println("requirements installed")
	}
	return nil
}

func (p *PythonLauncher) Run(ctx context.Context, job *scripts.Job) (scripts.Result, error) {
	targetDir := filepath.Dir(job.URL())

	err := p.installVenv(ctx, targetDir)
	if err != nil {
		_ = os.RemoveAll(targetDir)
		return scripts.Result{}, err
	}

	args := []string{filepath.Base(job.URL())}
	interpreter := filepath.Join(targetDir, "venv", "bin", "python")
	log.Println("Run interpreter", interpreter)

	values := job.Input()
	rawValues := make([]string, len(values))
	for i, v := range values {
		rawValues[i] = v.String()
	}
	args = append(args, rawValues...)

	var stdout, stderr bytes.Buffer
	var exitCode int

	cmd := exec.CommandContext(ctx, interpreter, args...)
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
