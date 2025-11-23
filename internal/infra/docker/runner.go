package docker

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

type Runner struct {
	cli *client.Client
	l   *slog.Logger
}

type RunResult struct {
	Status  int
	Message string
}

func NewRunner(l *slog.Logger) (*Runner, error) {
	if l == nil {
		return nil, errors.New("nil logger")
	}
	cli, err := client.New(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}
	return &Runner{cli, l}, nil
}

func MustNewRunner(l *slog.Logger) *Runner {
	r, err := NewRunner(l)
	if err != nil {
		panic(err)
	}
	return r
}

func (r *Runner) Build(ctx context.Context, path string, image string) error {
	l := r.l.With(
		slog.String("op", "docker.Runner.Build"),
		slog.String("image", image),
	)

	buildCtx, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open archive %q: %w", path, err)
	}

	l.Debug("Docker build started")
	res, err := r.cli.ImageBuild(ctx, buildCtx, client.ImageBuildOptions{
		Tags:       []string{image},
		Dockerfile: "Dockerfile",
	})
	if err != nil {
		return fmt.Errorf("failed to build image: %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	// Не знаю почему, но сборка не идёт, если не прочесть res.Body
	_, _ = io.ReadAll(res.Body)

	l.Debug("Docker build finished")
	return nil
}

func (r *Runner) Run(ctx context.Context, image string, input string) (RunResult, error) {
	l := r.l.With(
		slog.String("op", "docker.Runner.Run"),
		slog.String("image", image),
	)

	l.Debug("Docker container creating started")
	resp, err := r.cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Image: image,
		Config: &container.Config{
			OpenStdin:   true,
			AttachStdin: true,
			StdinOnce:   true,
		},
	})
	if err != nil {
		return RunResult{}, fmt.Errorf("failed to create container: %w", err)
	}
	l.Debug("Docker container created")

	l.Debug("Docker container starting")
	_, err = r.cli.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{})
	if err != nil {
		return RunResult{}, fmt.Errorf("failed to start container: %w", err)
	}
	l.Debug("Docker container started")

	l.Debug("Docker container attaching")
	attach, err := r.cli.ContainerAttach(ctx, resp.ID, client.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
	})
	if err != nil {
		return RunResult{}, fmt.Errorf("failed to attach container: %w", err)
	}
	l.Debug("Docker container attached")

	l.Debug("Docker container writing")
	n, err := attach.Conn.Write([]byte(input))
	if err != nil {
		_ = attach.Conn.Close()
		return RunResult{}, fmt.Errorf("failed to write input: %w", err)
	}
	l.Debug("Docker container input written", slog.Int("bytes", n))

	l.Debug("Docker container waiting")
	wRes := r.cli.ContainerWait(ctx, resp.ID, client.ContainerWaitOptions{
		Condition: container.WaitConditionNotRunning,
	})

	var result RunResult
	select {
	case err := <-wRes.Error:
		if err != nil {
			return RunResult{}, fmt.Errorf("failed to wait container: %w", err)
		}
	case res := <-wRes.Result:
		result.Status = int(res.StatusCode)
	}
	l.Debug("Docker container exited", slog.Int("status", result.Status))

	out, err := r.cli.ContainerLogs(ctx, resp.ID, client.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
	})
	if err != nil {
		return result, fmt.Errorf("failed to get container logs: %w", err)
	}
	defer func() { _ = out.Close() }()

	output, err := r.readDockerLogs(out)
	if err != nil {
		return result, fmt.Errorf("failed to get container logs: %w", err)
	}
	result.Message = output

	_, err = r.cli.ContainerRemove(ctx, resp.ID, client.ContainerRemoveOptions{})
	if err != nil {
		return result, fmt.Errorf("failed to remove container: %w", err)
	}

	return result, nil
}

func (r *Runner) readDockerLogs(rd io.Reader) (string, error) {
	br := bufio.NewReader(rd)
	var builder strings.Builder
	for {
		_, err := br.Discard(8)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", fmt.Errorf("failed to read Docker 8-byte header in output: %w", err)
		}
		line, err := br.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to readline Docker output: %w", err)
		}
		builder.WriteString(line)
	}
	return builder.String(), nil
}

func (r *Runner) Cleanup(ctx context.Context, image string) error {
	_, err := r.cli.ImageRemove(ctx, image, client.ImageRemoveOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove image: %w", err)
	}
	return nil
}
