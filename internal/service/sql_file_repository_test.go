package service_test

import (
	"os"
	"testing"

	"github.com/bmstu-itstech/scriptum-back/internal/service"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func setUpFileRepository() (*service.FileRepo, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, err
	}
	dsn := os.Getenv("DATABASE_URI")

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return service.NewFileRepository(db), err
}

func testFileRepository(t *testing.T) {
	r, err := setUpFileRepository()
	require.NoError(t, err)

	fileRepository_FileNotFound(t, r)
	fileRepository_Create(t, r)
	fileRepository_FileFound(t, r)
}
