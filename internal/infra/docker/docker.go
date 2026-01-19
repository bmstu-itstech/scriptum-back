package docker

import (
	"context"

	"github.com/moby/moby/client"
)

// TODO: add check in app

func IsDockerAvailable() bool {
	cli, err := client.New(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return false
	}
	_, err = cli.Ping(context.Background(), client.PingOptions{})
	return err == nil
}
