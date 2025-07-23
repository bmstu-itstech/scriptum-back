package service

import (
	"context"
	"testing"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"
)

func TestPostJob(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewConn()
	require.NoError(t, err)

	repo := &JobRepo{DB: mock}

	intVal, _ := scripts.NewInteger(42)
	inVec, _ := scripts.NewVector([]scripts.Value{intVal})

	jobID := scripts.JobID(0)
	userID := scripts.UserID(10)
	command := "test"
	startedAt := time.Now()
	scriptID := scripts.ScriptID(5)

	job, err := scripts.NewEmptyJob(jobID, userID, *inVec, command, startedAt)
	require.NoError(t, err)

	expectedID := scripts.JobID(123)

	mock.ExpectQuery("INSERT INTO jobs").
		WithArgs(job.UserID(), scriptID).
		WillReturnRows(pgxmock.NewRows([]string{"job_id"}).AddRow(int(expectedID)))

	id, err := repo.Post(ctx, *job, scriptID)
	require.NoError(t, err)
	require.Equal(t, expectedID, id)
}
func TestCloseJob(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewConn()
	require.NoError(t, err)

	repo := &JobRepo{DB: mock}

	jobID := scripts.JobID(123)
	userID := scripts.UserID(10)
	startedAt := time.Now()

	intVal, _ := scripts.NewInteger(1)
	inVec, _ := scripts.NewVector([]scripts.Value{intVal})

	realVal, _ := scripts.NewReal(3.14)
	outVec, _ := scripts.NewVector([]scripts.Value{realVal})

	job, err := scripts.NewEmptyJob(jobID, userID, *inVec, "cmd", startedAt)
	require.NoError(t, err)

	errorMsg := scripts.ErrorMessage("some error")
	status := scripts.StatusCode(1)

	result, err := scripts.NewResult(*job, status, *outVec, &errorMsg, time.Now())
	require.NoError(t, err)

	mock.ExpectBegin()

	mock.ExpectExec("UPDATE jobs SET").
		WithArgs(status, *result.ErrorMessage(), jobID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	mock.ExpectQuery("INSERT INTO parameters").
		WithArgs("3.14").
		WillReturnRows(pgxmock.NewRows([]string{"parameter_id"}).AddRow(int64(555)))

	mock.ExpectExec("INSERT INTO job_params").
		WithArgs(jobID, int64(555)).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	mock.ExpectCommit()

	err = repo.Update(ctx, jobID, result)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet())
}
