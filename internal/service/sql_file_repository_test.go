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

	return service.NewFileRepository(db, nil), err
}

func testFileRepository(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test sql file repository in short mode.")
	}
	r, err := setUpFileRepository()
	require.NoError(t, err)

	s, err := setUpScriptRepository()
	require.NoError(t, err)

	fileRepository_FileNotFound(t, r)
	fileRepository_Create(t, r)
	scriptRepository_Create(t, s)
	fileRepository_FileFound(t, r)
}
