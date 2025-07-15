package service

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/require"
)

func TestGetResult(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer mock.Close(ctx)

	repo := &JobRepo{DB: mock}

	jobID := scripts.JobID(2)
	startedAt := time.Now()
	userID := scripts.UserID(42)
	scriptID := int64(99)
	statusCode := 0
	errorMessage := ""

	valStr := "3.14"
	paramType := "out"
	fieldType := "real"
	valStrIn := "42"
	paramTypeIn := "in"
	fieldTypeIn := "integer"

	mock.ExpectQuery("SELECT(.*)FROM jobs").
		WithArgs(jobID).
		WillReturnRows(pgxmock.NewRows([]string{
			"user_id", "started_at", "status_code", "error_message", "script_id", "value", "param", "field_type",
		}).
			AddRow(userID, startedAt, statusCode, errorMessage, scriptID, &valStrIn, &paramTypeIn, fieldTypeIn).
			AddRow(userID, startedAt, statusCode, errorMessage, scriptID, &valStr, &paramType, fieldType),
		)

	res, err := repo.GetResult(ctx, jobID)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, userID, res.Job().UserID())
	require.Equal(t, scripts.StatusCode(statusCode), res.Code())
	require.Equal(t, errorMessage, *res.ErrorMessage())
	require.Equal(t, 1, res.Out().Len())
}

func TestGetResultsForUser(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer mock.Close(ctx)

	repo := &JobRepo{DB: mock}

	userID := scripts.UserID(42)
	startedAt := time.Now()
	statusCode := 0
	errorMessage := ""
	scriptID := int64(99)

	valStrIn := "42"
	paramTypeIn := "in"
	fieldTypeIn := "integer"
	valStrOut := "3.14"
	paramTypeOut := "out"
	fieldTypeOut := "real"

	rows := pgxmock.NewRows([]string{
		"job_id", "user_id", "started_at", "status_code", "error_message", "script_id", "value", "param", "field_type",
	}).
		AddRow(1, userID, startedAt, statusCode, errorMessage, scriptID, &valStrIn, &paramTypeIn, fieldTypeIn).
		AddRow(1, userID, startedAt, statusCode, errorMessage, scriptID, &valStrOut, &paramTypeOut, fieldTypeOut).
		AddRow(2, userID, startedAt, statusCode, errorMessage, scriptID, &valStrIn, &paramTypeIn, fieldTypeIn).
		AddRow(2, userID, startedAt, statusCode, errorMessage, scriptID, &valStrOut, &paramTypeOut, fieldTypeOut)

	mock.ExpectQuery("SELECT(.*)FROM jobs").
		WithArgs(userID).
		WillReturnRows(rows)

	results, err := repo.GetResultsForUser(ctx, userID)
	require.NoError(t, err)
	require.Len(t, results, 2)

	res1 := results[0]
	vec := res1.Job().In()
	require.Equal(t, userID, res1.Job().UserID())
	require.Equal(t, scripts.StatusCode(statusCode), res1.Code())
	require.Equal(t, errorMessage, *res1.ErrorMessage())
	require.Equal(t, 1, vec.Len())
	require.Equal(t, 1, res1.Out().Len())

	res2 := results[1]
	vec = res2.Job().In()
	require.Equal(t, userID, res2.Job().UserID())
	require.Equal(t, scripts.StatusCode(statusCode), res2.Code())
	require.Equal(t, errorMessage, *res2.ErrorMessage())
	require.Equal(t, 1, vec.Len())
	require.Equal(t, 1, res2.Out().Len())
}

func TestPostJob(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer mock.Close(ctx)

	repo := &JobRepo{DB: mock}

	intVal, _ := scripts.NewInteger(42)
	inVec, _ := scripts.NewVector([]scripts.Value{intVal})

	jobID := scripts.JobID(0)
	userID := scripts.UserID(10)
	command := "test"
	startedAt := time.Now()
	scriptID := scripts.ScriptID(5)

	job, err := scripts.NewJob(jobID, userID, *inVec, command, startedAt)
	fmt.Println(reflect.TypeOf(jobID))
	require.NoError(t, err)

	expectedID := scripts.JobID(123)

	mock.ExpectQuery("INSERT INTO jobs").
		WithArgs(job.UserID(), scriptID).
		WillReturnRows(pgxmock.NewRows([]string{"job_id"}).AddRow(int(expectedID)))

	id, err := repo.PostJob(ctx, *job, scriptID)
	require.NoError(t, err)
	require.Equal(t, expectedID, id)
}

func TestCloseJob(t *testing.T) {
	ctx := context.Background()
	mock, err := pgxmock.NewConn()
	require.NoError(t, err)
	defer mock.Close(ctx)

	repo := &JobRepo{DB: mock}

	jobID := scripts.JobID(123)
	userID := scripts.UserID(10)
	startedAt := time.Now()

	intVal, _ := scripts.NewInteger(1)
	inVec, _ := scripts.NewVector([]scripts.Value{intVal})

	realVal, _ := scripts.NewReal(3.14)
	outVec, _ := scripts.NewVector([]scripts.Value{realVal})

	job, err := scripts.NewJob(jobID, userID, *inVec, "cmd", startedAt)
	require.NoError(t, err)

	errorMsg := scripts.ErrorMessage("some error")
	status := scripts.StatusCode(1)

	result, err := scripts.NewResult(*job, status, *outVec, errorMsg)
	require.NoError(t, err)

	mock.ExpectExec("UPDATE jobs SET").
		WithArgs(status, *result.ErrorMessage(), jobID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.CloseJob(ctx, jobID, result)
	require.NoError(t, err)
}
