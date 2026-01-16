package docker_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/bmstu-itstech/scriptum-back/internal/config"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
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

	archivePath, err := testutils.TarCreate("tests/adder")
	require.NoError(t, err)
	t.Cleanup(func() {
		err = os.Remove(archivePath)
		if err != nil {
			t.Logf("failed to remove test archive: %v", err)
		}
	})
	buildCtx, err := os.Open(archivePath)
	require.NoError(t, err)

	l := slogdiscard.NewDiscardLogger()
	ctx := context.Background()
	ctx, cancelFn := context.WithTimeout(ctx, time.Second*10)
	defer cancelFn()

	cfg := config.Docker{
		ImagePrefix:   "sc-box",
		RunnerTimeout: 10 * time.Second,
	}
	r := docker.MustNewRunner(cfg, l)
	image, err := r.Build(ctx, buildCtx, value.NewBoxID())
	require.NoError(t, err)
	t.Cleanup(func() {
		err = r.Cleanup(context.Background(), image)
		if err != nil {
			t.Logf("failed to cleanup image: %v", err)
		}
	})

	t.Run("successfully added", func(t *testing.T) {
		res, err2 := r.Run(ctx, image, []value.Value{
			value.MustNewIntegerValue("1"),
			value.MustNewIntegerValue("2"),
		})
		require.NoError(t, err2)
		require.Equal(t, value.NewResult(0).WithOutput("3\n"), res)
	})

	t.Run("should return exception on invalid input", func(t *testing.T) {
		res, err2 := r.Run(ctx, image, []value.Value{
			value.MustNewIntegerValue("1"),
			value.NewStringValue("a"),
		})
		require.NoError(t, err2)
		require.NotEqual(t, value.ExitCode(0), res.Code())
		require.NotEmpty(t, res.Output())
	})

	t.Run("should return error if image not found", func(t *testing.T) {
		_, err = r.Run(ctx, "invalid", []value.Value{
			value.MustNewIntegerValue("1"),
			value.MustNewIntegerValue("2"),
		})
		require.Error(t, err)
	})
}
