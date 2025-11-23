package docker_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bmstu-itstech/scriptum-back/internal/infra/docker"
	"github.com/bmstu-itstech/scriptum-back/pkg/logs/handlers/slogdiscard"
	"github.com/bmstu-itstech/scriptum-back/pkg/testutils"
)

func TestRunner_Adder(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	if !docker.IsDockerAvailable() {
		t.Skip("docker is not available")
	}

	archive, err := testutils.TarCreate("tests/adder")
	require.NoError(t, err)
	t.Cleanup(func() {
		err = os.Remove(archive)
		if err != nil {
			t.Logf("failed to remove test archive: %v", err)
		}
	})

	l := slogdiscard.NewDiscardLogger()
	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, time.Second*10)
	defer cancelFn()

	r := docker.MustNewRunner(l)
	image := "sc-test-adder"
	err = r.Build(ctx, archive, image)
	require.NoError(t, err)
	t.Cleanup(func() {
		err = r.Cleanup(ctx, image)
		if err != nil {
			t.Logf("failed to cleanup image: %v", err)
		}
	})

	t.Run("successfully added", func(t *testing.T) {
		res, err := r.Run(ctx, image, "1\n2\n")
		require.NoError(t, err)
		require.Equal(t, docker.RunResult{
			Status:  0,
			Message: "3\n",
		}, res)
	})

	t.Run("should return exception on invalid input", func(t *testing.T) {
		res, err := r.Run(ctx, image, "1\na\n")
		require.NoError(t, err)
		require.NotEqual(t, 0, res.Status)
		require.NotEmpty(t, res.Message)
	})

	t.Run("should return error if image not found", func(t *testing.T) {
		_, err := r.Run(ctx, "", "1\n")
		require.Error(t, err)
	})
}
