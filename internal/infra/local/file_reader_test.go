package local_test

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/bmstu-itstech/scriptum-back/internal/config"
	"github.com/bmstu-itstech/scriptum-back/internal/domain/value"
	"github.com/bmstu-itstech/scriptum-back/internal/infra/local"
	"github.com/bmstu-itstech/scriptum-back/pkg/logs/handlers/slogdiscard"
)

func TestFileReader_Read(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	cfg := config.Storage{BasePath: "tests"}
	store := local.MustNewStorage(cfg, slogdiscard.NewDiscardLogger())
	fileID := value.FileID("1234abcd")

	rc, err := store.Read(t.Context(), fileID)
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, rc.Close())
	})

	bs, err := io.ReadAll(rc)
	require.NoError(t, err)

	s := string(bs)
	require.Equal(t, "1234abcd\n", s)
}
