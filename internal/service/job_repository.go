package service

import (
	"context"
	"fmt"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

type JobDBConn interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

type JobRepo struct {
	DB JobDBConn
}

func NewJobRepo(ctx context.Context) (*JobRepo, error) {
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
	return &JobRepo{
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

func (r *JobRepo) GetResult(ctx context.Context, jobID scripts.JobID) (scripts.Result, error) {
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

const PostJobQuery = `
		INSERT INTO jobs (user_id, script_id, started_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		RETURNING job_id;
	`

func (r *JobRepo) PostJob(ctx context.Context, job scripts.Job, scriptID scripts.ScriptID) (scripts.JobID, error) {
	var rawID int
	err := r.DB.QueryRow(ctx, PostJobQuery,
		job.UserID(),
		scriptID,
	).Scan(&rawID)
	if err != nil {
		return 0, err
	}
	return scripts.JobID(rawID), nil
}

const CloseJobQuery = `
	UPDATE jobs SET
		status_code = $1,
		error_message = $2,
		closed_at = CURRENT_TIMESTAMP
	WHERE job_id = $3;
`

const insertOutParamQuery = `
		INSERT INTO parameters (value) VALUES ($1) RETURNING parameter_id;
	`

const insertJobParamQuery = `
		INSERT INTO job_params (job_id, parameter_id, param)
		VALUES ($1, $2, 'out');
	`

func (r *JobRepo) CloseJob(ctx context.Context, jobID scripts.JobID, res *scripts.Result) error {
	tx, err := r.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, CloseJobQuery,
		res.Code(),
		*res.ErrorMessage(),
		jobID,
	)

	if err != nil {
		return err
	}

	for _, val := range res.Out().Get() {
		var paramID int64
		err := tx.QueryRow(ctx, insertOutParamQuery, val).Scan(&paramID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, insertJobParamQuery, jobID, paramID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

const JobsByScriptIDQuery = `
	SELECT
		j.job_id,
		j.user_id,
		j.started_at,
		j.script_id,
		p.value,
		jp.param,
		f.field_type
	FROM jobs j
	LEFT JOIN job_params jp ON j.job_id = jp.job_id
	LEFT JOIN parameters p ON p.parameter_id = jp.parameter_id
	LEFT JOIN fields f ON f.field_id = p.field_id
	WHERE j.script_id = $1;
`

type JobAccumulator struct {
	userID    scripts.UserID
	startedAt time.Time
	inputVals []scripts.Value
}

func (r *JobRepo) JobsByScriptID(ctx context.Context, scriptID scripts.ScriptID) ([]scripts.Job, error) {
	rows, err := r.DB.Query(ctx, JobsByScriptIDQuery, scriptID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	jobMap := make(map[scripts.JobID]*JobAccumulator)

	for rows.Next() {
		var (
			jobID     int
			userID    scripts.UserID
			startedAt time.Time

			valStr    *string
			paramType *string
			fieldType string
		)

		if err := rows.Scan(&jobID, &userID, &startedAt, new(int64), &valStr, &paramType, &fieldType); err != nil {
			return nil, err
		}

		jid := scripts.JobID(jobID)
		acc, ok := jobMap[jid]
		if !ok {
			acc = &JobAccumulator{
				userID:    userID,
				startedAt: startedAt,
			}
			jobMap[jid] = acc
		}

		if valStr != nil && paramType != nil && *paramType == "in" {
			val, err := scripts.NewValue(fieldType, *valStr)
			if err != nil {
				return nil, err
			}
			acc.inputVals = append(acc.inputVals, val)
		}
	}

	var jobs []scripts.Job
	for jid, acc := range jobMap {
		inVec, err := scripts.NewVector(acc.inputVals)
		if err != nil {
			return nil, err
		}
		job, err := scripts.NewJob(jid, acc.userID, *inVec, "", acc.startedAt)
		if err != nil {
			return nil, err
		}
		jobs = append(jobs, *job)
	}

	return jobs, nil
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

func (r *JobRepo) GetResultsForUser(ctx context.Context, userID scripts.UserID) ([]scripts.Result, error) {
	return r.getResultsBase(ctx, GetResultsForUserQuery, userID)
}

func (r *JobRepo) SearchJobs(ctx context.Context, userID scripts.UserID, substr string) ([]scripts.Result, error) {
	return r.getResultsBase(ctx, SearchJobsQuery, userID, substr)
}
func (r *JobRepo) getResultsBase(ctx context.Context, query string, args ...any) ([]scripts.Result, error) {
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
		default:
			return nil, fmt.Errorf("unknown query")
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
