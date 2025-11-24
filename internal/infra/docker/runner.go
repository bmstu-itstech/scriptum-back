package docker

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
)

const imagePrefix = "sc-box"

type Runner struct {
	cli *client.Client
	l   *slog.Logger
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

func (r *Runner) Build(ctx context.Context, buildCtx io.Reader, id value.BoxID) (value.ImageTag, error) {
	l := r.l.With(
		slog.String("op", "docker.Runner.Build"),
		slog.String("box_id", string(id)),
	)

	image := value.NewImageTag(imagePrefix, id)
	l = l.With(slog.String("image", string(image)))

	l.Debug("Docker build started")
	res, err := r.cli.ImageBuild(ctx, buildCtx, client.ImageBuildOptions{
		Tags:       []string{string(image)},
		Dockerfile: "Dockerfile",
	})
	if err != nil {
		return "", fmt.Errorf("failed to build image: %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	// Не знаю почему, но сборка не идёт, если не прочесть res.Body
	_, _ = io.ReadAll(res.Body)

	l.Debug("Docker build finished")
	return image, nil
}

func (r *Runner) Run(ctx context.Context, image value.ImageTag, input []value.Value) (value.Result, error) {
	l := r.l.With(
		slog.String("op", "docker.Runner.Run"),
		slog.String("image", string(image)),
	)

	l.Debug("Docker container creating started")
	resp, err := r.cli.ContainerCreate(ctx, client.ContainerCreateOptions{
		Image: string(image),
		Config: &container.Config{
			OpenStdin:   true,
			AttachStdin: true,
			StdinOnce:   true,
		},
	})
	if err != nil {
		return value.Result{}, fmt.Errorf("failed to create container: %w", err)
	}
	l.Debug("Docker container created")

	l.Debug("Docker container starting")
	_, err = r.cli.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{})
	if err != nil {
		return value.Result{}, fmt.Errorf("failed to start container: %w", err)
	}
	l.Debug("Docker container started")

	l.Debug("Docker container attaching")
	attach, err := r.cli.ContainerAttach(ctx, resp.ID, client.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
	})
	if err != nil {
		return value.Result{}, fmt.Errorf("failed to attach container: %w", err)
	}
	l.Debug("Docker container attached")

	l.Debug("Docker container writing")
	n, err := attach.Conn.Write(r.marshallInput(input))
	if err != nil {
		_ = attach.Conn.Close()
		return value.Result{}, fmt.Errorf("failed to write input: %w", err)
	}
	l.Debug("Docker container input written", slog.Int("bytes", n))

	l.Debug("Docker container waiting")
	wRes := r.cli.ContainerWait(ctx, resp.ID, client.ContainerWaitOptions{
		Condition: container.WaitConditionNotRunning,
	})

	var result value.Result
	select {
	case err := <-wRes.Error:
		if err != nil {
			return value.Result{}, fmt.Errorf("failed to wait container: %w", err)
		}
	case res := <-wRes.Result:
		result = value.NewResult(value.ExitCode(res.StatusCode))
	}
	l.Debug("Docker container exited", slog.Int("exit_code", int(result.Code())))

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
	result = result.WithOutput(output)

	_, err = r.cli.ContainerRemove(ctx, resp.ID, client.ContainerRemoveOptions{})
	if err != nil {
		return result, fmt.Errorf("failed to remove container: %w", err)
	}

	return result, nil
}

func (r *Runner) marshallInput(input []value.Value) []byte {
	var buf bytes.Buffer
	for _, v := range input {
		buf.WriteString(v.String())
		buf.WriteRune('\n')
	}
	return buf.Bytes()
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

func (r *Runner) Cleanup(ctx context.Context, image value.ImageTag) error {
	_, err := r.cli.ImageRemove(ctx, string(image), client.ImageRemoveOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove image: %w", err)
	}
	return nil
}
