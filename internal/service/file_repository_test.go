package service_test

import (
	"context"
	"testing"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/stretchr/testify/require"
)

func fileRepository_Create(t *testing.T, repo scripts.FileRepository) {
	ctx := context.Background()

	url := scripts.URL("http://example.com/file1")
	id, err := repo.Create(ctx, &url)
	require.NoError(t, err)
	require.NotZero(t, id)

	url = scripts.URL("http://example.com/file2")
	id, err = repo.Create(ctx, &url)
	require.NoError(t, err)
	require.NotZero(t, id)

	url = scripts.URL("http://example.com/file3")
	id, err = repo.Create(ctx, &url)
	require.NoError(t, err)
	require.NotZero(t, id)
}

func fileRepository_FileFound(t *testing.T, repo scripts.FileRepository) {
	ctx := context.Background()

	var existingFileID scripts.FileID = 1

	file, err := repo.File(ctx, existingFileID)
	require.NoError(t, err)
	require.NotNil(t, file)
	require.Equal(t, existingFileID, file.ID())
}

func fileRepository_FileNotFound(t *testing.T, repo scripts.FileRepository) {
	ctx := context.Background()

	nonExistentFileID := scripts.FileID(999999999)

	file, err := repo.File(ctx, nonExistentFileID)
	require.Error(t, err)
	require.Nil(t, file)
	require.ErrorIs(t, err, scripts.ErrFileNotFound)
}
