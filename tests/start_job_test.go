//nolint:testpackage // именно такое название и нужно
package tests

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	apiv2 "github.com/bmstu-itstech/scriptum-back/gen/go/api/v2"
	"github.com/bmstu-itstech/scriptum-back/pkg/testutils"
	tsuite "github.com/bmstu-itstech/scriptum-back/tests/suite"
)

func createAdderBox(ctx context.Context, suite *tsuite.Suite) (string, error) {
	archivePath, err := testutils.TarCreate("res/adder")
	if err != nil {
		return "", err
	}
	filename := filepath.Base(archivePath)
	defer func() {
		_ = os.Remove(archivePath)
	}()

	f, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}

	archiveID, err := UploadFile(ctx, suite.FileService, filename, f)
	if err != nil {
		return "", err
	}

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
	if err != nil {
		return "", err
	}

	return res.BoxId, nil
}

func TestJobServiceStartJob(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	ctx, suite := tsuite.New(t)

	boxID, err := createAdderBox(ctx, suite)
	require.NoError(t, err)

	res, err := suite.BoxService.StartJob(ctx, &apiv2.StartJobRequest{
		BoxId: boxID,
		Values: []*apiv2.Value{
			{
				Type:  apiv2.Type_TYPE_INTEGER,
				Value: "1",
			},
			{
				Type:  apiv2.Type_TYPE_INTEGER,
				Value: "2",
			},
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, res.GetJobId())

	require.Eventually(t, func() bool {
		res2, err2 := suite.JobService.GetJob(ctx, &apiv2.GetJobRequest{JobId: res.GetJobId()})
		require.NoError(t, err2)
		require.NotNil(t, res2)
		// Есть возможность, что job не закончила работать.
		// Однако, если job закончилась, то необходимо проверить инварианты через require, так как повторять запрос
		// в случае, если они нарушены, не имеет смысла.
		if res2.Job.State == apiv2.JobState_JOB_STATE_FINISHED {
			require.NotNil(t, res2.Job.Result)
			jRes := res2.Job.Result
			require.Equal(t, int64(0), jRes.Code)
			require.Len(t, jRes.Output, 1)
			require.Equal(t, apiv2.Type_TYPE_INTEGER, jRes.Output[0].Type)
			require.Equal(t, "3", jRes.Output[0].Value)
			return true
		}
		return false
	}, 5*time.Second, time.Second)
}
