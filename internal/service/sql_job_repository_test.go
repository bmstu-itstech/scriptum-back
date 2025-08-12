package service_test

import (
	"os"
	"testing"

	"github.com/bmstu-itstech/scriptum-back/internal/service"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func setUpJobRepository() (*service.JobRepo, error) {
	err := godotenv.Load("../../.env")
	if err != nil {
		return nil, err
	}
	dsn := os.Getenv("DATABASE_URI")

	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return service.NewJobRepository(db), err
}

func testJobRepository(t *testing.T) {
	r, err := setUpJobRepository()
	require.NoError(t, err)
	jobRepository_JobNotFound(t, r)
	jobRepository_UserJobs_NotFound(t, r)
	jobRepository_UserJobsWithState_NotFound(t, r)

	jobRepository_JobFound(t, r)
	jobRepository_UserJobsWithState_Found(t, r)
	jobRepository_UserJobs_Found(t, r)

	jobRepository_Create(t, r)
	jobRepository_CreateMultiple(t, r)

	jobRepository_Update(t, r)

	jobRepository_Delete(t, r)

	jobRepository_MixedUserJobs(t, r)

	jobRepository_MixedUserJobsWithState(t, r)

}
