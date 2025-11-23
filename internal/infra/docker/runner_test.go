package docker_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

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
	image, err := r.Build(ctx, archive, value.NewBoxID())
	require.NoError(t, err)
	t.Cleanup(func() {
		err = r.Cleanup(ctx, image)
		if err != nil {
			t.Logf("failed to cleanup image: %v", err)
		}
	})

	t.Run("successfully added", func(t *testing.T) {
		input := value.NewEmptyInput().With(value.MustNewIntegerValue("1")).With(value.MustNewIntegerValue("2"))
		res, err := r.Run(ctx, image, input)
		require.NoError(t, err)
		require.Equal(t, value.NewResult(0).WithOutput("3\n"), res)
	})

	t.Run("should return exception on invalid input", func(t *testing.T) {
		input := value.NewEmptyInput().With(value.MustNewIntegerValue("1")).With(value.NewStringValue("a"))
		res, err := r.Run(ctx, image, input)
		require.NoError(t, err)
		require.NotEqual(t, value.ExitCode(0), res.Code())
		require.NotEmpty(t, res.Output())
	})

	t.Run("should return error if image not found", func(t *testing.T) {
		input := value.NewEmptyInput().With(value.MustNewIntegerValue("1")).With(value.MustNewIntegerValue("2"))
		_, err := r.Run(ctx, "", input)
		require.Error(t, err)
	})
}
