//nolint:testpackage // именно такое название и нужно
package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	apiv2 "github.com/bmstu-itstech/scriptum-back/gen/go/api/v2"
	"github.com/bmstu-itstech/scriptum-back/pkg/testutils"
	tsuite "github.com/bmstu-itstech/scriptum-back/tests/suite"
)

func TestBoxServiceCreateBox(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	ctx, suite := tsuite.New(t)

	archivePath, err := testutils.TarCreate("res/adder")
	require.NoError(t, err)
	filename := filepath.Base(archivePath)
	t.Cleanup(func() {
		err = os.Remove(archivePath)
		if err != nil {
			t.Logf("failed to remove test archive: %v", err)
		}
	})

	f, err := os.Open(archivePath)
	require.NoError(t, err)

	archiveID, err := UploadFile(ctx, suite.FileService, filename, f)
	require.NoError(t, err)

	res, err := suite.BoxService.CreateBox(ctx, &apiv2.CreateBoxRequest{
		ArchiveId: archiveID,
		Name:      "Adder",
		Input: []*apiv2.Field{
			{
				Type: apiv2.Type_TYPE_INTEGER,
				Name: "A",
			},
			{
				Type: apiv2.Type_TYPE_INTEGER,
				Name: "B",
			},
		},
		Output: []*apiv2.Field{
			{
				Type: apiv2.Type_TYPE_INTEGER,
				Name: "A + B",
			},
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, res.GetBoxId())

	box, err := suite.BoxService.GetBox(ctx, &apiv2.GetBoxRequest{BoxId: res.GetBoxId()})
	require.NoError(t, err)
	require.NotEmpty(t, box)
	require.Equal(t, box.Box.GetId(), res.GetBoxId())
}
