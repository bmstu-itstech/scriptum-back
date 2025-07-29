package service_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/bmstu-itstech/scriptum-back/internal/service"
	"github.com/stretchr/testify/require"
)

func TestSystemManager_Save(t *testing.T) {
	// Создаём временную директорию для теста
	dir := "./tmp_test_dir"

	sm, err := service.NewSystemManager(dir)
	require.NoError(t, err)

	content := []byte("hello world")

	file, err := scripts.NewFile("testfile.txt", content)
	require.NoError(t, err)

	url, err := sm.Save(context.Background(), file)
	require.NoError(t, err)
	require.NotEmpty(t, url)

	// Проверяем, что файл создался и содержимое совпадает
	path := string(url)
	data, err := os.ReadFile(path)
	require.NoError(t, err)
	require.Equal(t, content, data)

	// Проверяем, что имя файла в нужной папке и содержит имя оригинала
	require.Contains(t, path, dir)
	require.Contains(t, path, "testfile.txt")

	// Дополнительно — файл можно удалить после теста, но TempDir уже чистит

	time.Sleep(1 * time.Second)
	err = os.Remove(path)
	require.NoError(t, err)
}
