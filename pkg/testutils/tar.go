package testutils

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"
)

// TarCreate создаёт архив .tar из всех файлов в директории и сохраняет его там же с названием директории.
func TarCreate(dirPath string) (string, error) {
	name := filepath.Base(dirPath) + ".tar"
	path := filepath.Join(dirPath, name)

	tarFile, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = tarFile.Close() }()

	tarWriter := tar.NewWriter(tarFile)
	defer func() { _ = tarWriter.Close() }()

	err = filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Пропускаем сам архив, чтобы не заархивировать его же самого
		if filePath == path {
			return nil
		}

		header, err := tar.FileInfoHeader(info, info.Name())
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(dirPath, filePath)
		if err != nil {
			return err
		}
		header.Name = relPath

		if err2 := tarWriter.WriteHeader(header); err2 != nil {
			return err2
		}

		if info.Mode().IsRegular() {
			file, err2 := os.Open(filePath)
			if err2 != nil {
				return err2
			}
			defer func() { _ = file.Close() }()

			_, err2 = io.Copy(tarWriter, file)
			if err2 != nil {
				return err2
			}
		}
		// Директории требуют только записи заголовка, поэтому здесь рассматриваются только регулярные файлы.

		return nil
	})

	if err != nil {
		return "", err
	}

	return path, nil
}
