package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/jackc/pgx/v4"
)

type ResRepo struct {
	DB SQLDBConn
}

func NewResRepo(ctx context.Context) (*ResRepo, error) {
	host := "localhost"
	port := 5432
	user := "app_user"
	password := "your_secure_password"
	dbname := "dev"

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, password, host, port, dbname)

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		return nil, err
	}
	return &ResRepo{
		DB: conn,
	}, nil
}

const GetResultQuery = `
	SELECT
		j.user_id,
		j.started_at,
		j.closed_at,
		j.status_code,
		j.error_message,
		j.script_id,
		p.value,
		jp.param,
		f.field_type
	FROM jobs j
	LEFT JOIN job_params jp ON j.job_id = jp.job_id
	LEFT JOIN parameters p ON p.parameter_id = jp.parameter_id
	LEFT JOIN fields f ON f.field_id = p.field_id
	WHERE j.job_id = $1;
`

func (r *ResRepo) GetResult(ctx context.Context, jobID scripts.JobID) (scripts.Result, error) {
	rows, err := r.DB.Query(ctx, GetResultQuery, jobID)
	if err != nil {
		return scripts.Result{}, err
	}
	defer rows.Close()

	var (
		userID       scripts.UserID
		startedAt    time.Time
		closed_at    time.Time
		statusCode   int
		errorMessage string
		scriptID     int64

		inputVals, outputVals []scripts.Value
	)

	for rows.Next() {
		var (
			valStr    *string
			paramType *string
			fieldType string
		)

		if err := rows.Scan(&userID, &startedAt, &closed_at, &statusCode, &errorMessage, &scriptID, &valStr, &paramType, &fieldType); err != nil {
			return scripts.Result{}, err
		}

		if valStr != nil && paramType != nil {
			val, err := scripts.NewValue(fieldType, *valStr)
			if err != nil {
				return scripts.Result{}, err
			}
			switch *paramType {
			case "in":
				inputVals = append(inputVals, val)
			case "out":
				outputVals = append(outputVals, val)
			}
		}
	}

	inVec, err := scripts.NewVector(inputVals)
	if err != nil {
		return scripts.Result{}, err
	}
	outVec, err := scripts.NewVector(outputVals)
	if err != nil {
		return scripts.Result{}, err
	}

	job, err := scripts.NewJob(jobID, userID, *inVec, "", startedAt)
	if err != nil {
		return scripts.Result{}, err
	}

	errMsg := scripts.ErrorMessage(errorMessage)

	result, err := scripts.NewResult(*job, scripts.StatusCode(statusCode), *outVec, &errMsg, closed_at)
	if err != nil {
		return scripts.Result{}, err
	}

	return *result, nil
}

const SearchJobsQuery = `
SELECT 
    j.job_id,
    j.started_at,
    j.closed_at,
    j.status_code,
    j.error_message,
    s.script_id,
    s.name,
    s.description,
    s.path,
    s.visibility,
    s.owner_id,
    s.created_at
FROM jobs j
JOIN scripts s ON s.script_id = j.script_id
WHERE j.user_id = $1
  AND s.name ILIKE '%' || $2 || '%'
ORDER BY j.started_at DESC;
`

const GetResultsForUserQuery = `
	SELECT
		j.job_id,
		j.user_id,
		j.started_at,
		j.status_code,
		j.error_message,
		j.script_id,
		p.value,
		jp.param,
		f.field_type
	FROM jobs j
	LEFT JOIN job_params jp ON j.job_id = jp.job_id
	LEFT JOIN parameters p ON p.parameter_id = jp.parameter_id
	LEFT JOIN fields f ON f.field_id = p.field_id
	WHERE j.user_id = $1
	ORDER BY j.started_at DESC;
`

type resultAccumulator struct {
	userID       scripts.UserID
	startedAt    time.Time
	closedAt     time.Time
	statusCode   int
	errorMessage *string
	scriptID     int64

	inputVals  []scripts.Value
	outputVals []scripts.Value
}

func (r *ResRepo) UserResults(ctx context.Context, userID scripts.UserID) ([]scripts.Result, error) {
	return r.getResultsBase(ctx, GetResultsForUserQuery, userID)
}

func (r *ResRepo) SearchResult(ctx context.Context, userID scripts.UserID, substr string) ([]scripts.Result, error) {
	return r.getResultsBase(ctx, SearchJobsQuery, userID, substr)
}
func (r *ResRepo) getResultsBase(ctx context.Context, query string, args ...any) ([]scripts.Result, error) {
	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobMap := make(map[scripts.JobID]*resultAccumulator)

	for rows.Next() {
		var (
			rawID       int
			userID      scripts.UserID
			startedAt   time.Time
			closedAt    time.Time
			statusCode  int
			errorMsg    string
			scriptID    int64
			scriptName  string
			scriptDesc  string
			scriptPath  string
			scriptVis   string
			scriptOwner int
			scriptDate  *time.Time

			valStr    *string
			paramType *string
			fieldType string
		)

		switch query {
		case SearchJobsQuery:
			err = rows.Scan(
				&rawID, &startedAt, &closedAt, &statusCode, &errorMsg,
				&scriptID, &scriptName, &scriptDesc, &scriptPath,
				&scriptVis, &scriptOwner, &scriptDate,
			)
		case GetResultsForUserQuery:
			err = rows.Scan(
				&rawID, &userID, &startedAt, &statusCode, &errorMsg,
				&scriptID, &valStr, &paramType, &fieldType,
			)
		}

		if err != nil {
			return nil, err
		}

		jobID := scripts.JobID(rawID)

		acc, ok := jobMap[jobID]
		if !ok {
			acc = &resultAccumulator{
				userID:       userID,
				startedAt:    startedAt,
				closedAt:     closedAt,
				statusCode:   statusCode,
				errorMessage: &errorMsg,
				scriptID:     scriptID,
			}
			jobMap[jobID] = acc
		}

		if valStr != nil && paramType != nil {
			val, err := scripts.NewValue(fieldType, *valStr)
			if err != nil {
				return nil, err
			}
			switch *paramType {
			case "in":
				acc.inputVals = append(acc.inputVals, val)
			case "out":
				acc.outputVals = append(acc.outputVals, val)
			}
		}
	}

	var results []scripts.Result
	for jobID, acc := range jobMap {
		inVec, err := scripts.NewVector(acc.inputVals)
		if err != nil {
			return nil, err
		}
		outVec, err := scripts.NewVector(acc.outputVals)
		if err != nil {
			return nil, err
		}

		job, err := scripts.NewJob(jobID, acc.userID, *inVec, "", acc.startedAt)
		if err != nil {
			return nil, err
		}

		errMsg := scripts.ErrorMessage(*acc.errorMessage)

		res, err := scripts.NewResult(*job, scripts.StatusCode(acc.statusCode), *outVec, &errMsg, acc.closedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, *res)
	}

	return results, nil
}
