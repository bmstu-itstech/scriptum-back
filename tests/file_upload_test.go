//nolint:testpackage // именно такое название и нужно
package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	apiv2 "github.com/bmstu-itstech/scriptum-back/gen/go/api/v2"
	tsuite "github.com/bmstu-itstech/scriptum-back/tests/suite"
)

func TestFileServiceUpload(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	ctx, suite := tsuite.New(t)

	s := "1234567890qwertyuiop\n"
	filename := "test.txt"

	stream, err := suite.FileService.Upload(ctx)
	require.NoError(t, err)

	err = stream.Send(&apiv2.FileUploadRequest{
		Body: &apiv2.FileUploadRequest_Meta{
			Meta: &apiv2.FileMeta{
				Name: filename,
			},
		},
	})
	require.NoError(t, err)

	data := []byte(s)
	testChunkSize := 8 // Для тестов проверим отправку чанками ОЧЕНЬ маленького размера (8 байт)

	for i := 0; i < len(data); i += testChunkSize {
		end := i + testChunkSize
		if end > len(data) {
			end = len(data)
		}

		err = stream.Send(&apiv2.FileUploadRequest{
			Body: &apiv2.FileUploadRequest_Chunk{
				Chunk: data[i:end],
			},
		})
		require.NoError(t, err)
	}

	res, err := stream.CloseAndRecv()
	require.NoError(t, err)

	require.NotEmpty(t, res.GetFileId())
	require.Equal(t, len(data), int(res.GetSize()))
}
