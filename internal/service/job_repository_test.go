package service_test

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/require"
)

func jobRepository_JobFound(t *testing.T, repo scripts.JobRepository) {
	proto := generateRandomJobPrototype(t)
	created, err := repo.Create(context.Background(), proto)
	require.NoError(t, err)

	got, err := repo.Job(context.Background(), created.ID())
	require.NoError(t, err)
	require.Equal(t, created.ID(), got.ID())
	require.Equal(t, created.OwnerID(), got.OwnerID())
	require.Equal(t, created.Input(), got.Input())
	require.Equal(t, created.State(), got.State())

}

func jobRepository_JobNotFound(t *testing.T, repo scripts.JobRepository) {
	_, err := repo.Job(context.Background(), 999999)
	require.Error(t, err)
}

func jobRepository_Create(t *testing.T, repo scripts.JobRepository) {
	proto := generateRandomJobPrototype(t)
	created, err := repo.Create(context.Background(), proto)
	require.NoError(t, err)
	require.NotNil(t, created)

	require.Equal(t, created.OwnerID(), proto.OwnerID())
	require.Equal(t, created.Input(), proto.Input())
}

func jobRepository_CreateMultiple(t *testing.T, repo scripts.JobRepository) {
	count := 5
	var prevID scripts.JobID = 0

	for i := 0; i < count; i++ {
		proto := generateRandomJobPrototype(t)
		created, err := repo.Create(context.Background(), proto)
		require.NoError(t, err)
		require.NotNil(t, created)

		require.Equal(t, created.OwnerID(), proto.OwnerID())
		require.Equal(t, created.Input(), proto.Input())

		require.Greater(t, created.ID(), prevID)
		prevID = created.ID()
	}
}

func jobRepository_Update(t *testing.T, repo scripts.JobRepository) {
	proto := generateRandomJobPrototype(t)
	created, err := repo.Create(context.Background(), proto)
	require.NoError(t, err)

	newState := "finished"

	output := generateRandomValuePrototype(t)

	res, err := scripts.NewSuccessResult(output)
	require.NoError(t, err)

	time := time.Now()

	val, err := scripts.NewValueType("integer")
	require.NoError(t, err)

	outputP1, err := scripts.NewField(*val, gofakeit.LetterN(10), gofakeit.LetterN(10), gofakeit.LetterN(10))
	require.NoError(t, err)

	updatedJob, err := scripts.RestoreJob(
		int64(created.ID()),
		int64(created.OwnerID()),
		int64(created.ScriptID()),
		newState,
		created.Input(),
		[]scripts.Field{*outputP1},
		"example",
		res,
		created.CreatedAt(),
		&time,
		nil,
	)
	require.NoError(t, err)

	err = repo.Update(context.Background(), updatedJob)
	require.NoError(t, err)

	got, err := repo.Job(context.Background(), created.ID())
	require.NoError(t, err)

	state, err := scripts.NewJobStateFromString("finished")
	require.NoError(t, err)

	gotRes, err := got.Result()
	require.NoError(t, err)

	require.Equal(t, created.ID(), got.ID())
	require.Equal(t, created.OwnerID(), got.OwnerID())
	require.Equal(t, created.Input(), got.Input())
	require.Equal(t, state, got.State())
	require.Equal(t, res, gotRes)
}

func jobRepository_Delete(t *testing.T, repo scripts.JobRepository) {
	proto := generateRandomJobPrototype(t)
	created, err := repo.Create(context.Background(), proto)
	require.NoError(t, err)

	err = repo.Delete(context.Background(), created.ID())
	require.NoError(t, err)

	_, err = repo.Job(context.Background(), created.ID())
	require.Error(t, err)
}

func jobRepository_UserJobs_Found(t *testing.T, repo scripts.JobRepository) {
	ctx := context.Background()
	ownerID := scripts.UserID(gofakeit.IntRange(1, 1000))
	scriptID := scripts.ScriptID(gofakeit.IntRange(1, 3))

	for i := 0; i < 2; i++ {
		input := []scripts.Value{}
		for j := 0; j < 2; j++ {
			v, err := scripts.NewValue("integer", strconv.Itoa(gofakeit.Number(0, 100)))
			require.NoError(t, err)
			input = append(input, v)
		}

		val, err := scripts.NewValueType("integer")
		require.NoError(t, err)

		outputP1, err := scripts.NewField(*val, gofakeit.LetterN(10), gofakeit.LetterN(10), gofakeit.LetterN(10))
		require.NoError(t, err)

		proto, err := scripts.NewJobPrototype(ownerID, scriptID, input, []scripts.Field{*outputP1}, "example", "")
		require.NoError(t, err)

		_, err = repo.Create(ctx, proto)
		require.NoError(t, err)
	}

	jobs, err := repo.UserJobs(ctx, ownerID)
	require.NoError(t, err)
	require.Len(t, jobs, 2)
	for _, j := range jobs {
		require.Equal(t, ownerID, j.OwnerID())
	}
}

func jobRepository_UserJobs_NotFound(t *testing.T, repo scripts.JobRepository) {
	ctx := context.Background()
	ownerID := scripts.UserID(999999)

	scriptsFound, err := repo.UserJobs(ctx, ownerID)
	require.NoError(t, err)
	require.Nil(t, scriptsFound)
}

func jobRepository_UserJobsWithState_Found(t *testing.T, repo scripts.JobRepository) {
	ctx := context.Background()
	ownerID := scripts.UserID(gofakeit.IntRange(1, 1000))

	for i := 0; i < 4; i++ {
		proto := generateRandomJobPrototype(t)
		protoOwner := ownerID
		if i%2 == 0 {
			protoOwner = scripts.UserID(ownerID + 1)
		}
		proto, err := scripts.NewJobPrototype(protoOwner, proto.ScriptID(), proto.Input(), proto.Expected(), proto.URL(), "")
		require.NoError(t, err)

		job, err := repo.Create(ctx, proto)
		require.NoError(t, err)

		if job.OwnerID() == ownerID {
			err := job.Run()
			require.NoError(t, err)
			err = repo.Update(ctx, job)
			require.NoError(t, err)
		}

	}

	jobs, err := repo.UserJobsWithState(ctx, ownerID, scripts.JobRunning)
	require.NoError(t, err)
	require.NotNil(t, jobs)
	require.NotEmpty(t, jobs)

	for _, job := range jobs {
		require.Equal(t, ownerID, job.OwnerID())
		require.Equal(t, scripts.JobRunning, job.State())
	}
}

func jobRepository_UserJobsWithState_NotFound(t *testing.T, repo scripts.JobRepository) {
	ctx := context.Background()
	ownerID := scripts.UserID(gofakeit.IntRange(1, 1000))

	for i := 0; i < 3; i++ {
		proto := generateRandomJobPrototype(t)
		proto, err := scripts.NewJobPrototype(ownerID, proto.ScriptID(), proto.Input(), proto.Expected(), proto.URL(), "")
		require.NoError(t, err)

		_, err = repo.Create(ctx, proto)
		require.NoError(t, err)
	}

	jobs, err := repo.UserJobsWithState(ctx, ownerID, scripts.JobRunning)
	require.NoError(t, err)
	require.Nil(t, jobs)
}

func jobRepository_MixedUserJobs(t *testing.T, repo scripts.JobRepository) {
	ctx := context.Background()
	userA := scripts.UserID(1)
	userB := scripts.UserID(2)

	for i := 0; i < 6; i++ {
		owner := userA
		if i%2 == 0 {
			owner = userB
		}

		proto := generateRandomJobPrototype(t)
		proto, err := scripts.NewJobPrototype(owner, proto.ScriptID(), proto.Input(), proto.Expected(), proto.URL(), "")
		require.NoError(t, err)

		_, err = repo.Create(ctx, proto)
		require.NoError(t, err)
	}

	jobsA, err := repo.UserJobs(ctx, userA)
	require.NoError(t, err)
	require.NotNil(t, jobsA)
	require.Len(t, jobsA, 3)

	for _, job := range jobsA {
		require.Equal(t, userA, job.OwnerID())
	}

	jobsB, err := repo.UserJobs(ctx, userB)
	require.NoError(t, err)
	require.NotNil(t, jobsB)
	require.Len(t, jobsB, 3)

	for _, job := range jobsB {
		require.Equal(t, userB, job.OwnerID())
	}
}

func jobRepository_MixedUserJobsWithState(t *testing.T, repo scripts.JobRepository) {
	ctx := context.Background()
	userA := scripts.UserID(10)
	userB := scripts.UserID(20)

	for i := 0; i < 6; i++ {
		owner := userA
		if i%2 == 0 {
			owner = userB
		}

		proto := generateRandomJobPrototype(t)
		proto, err := scripts.NewJobPrototype(owner, proto.ScriptID(), proto.Input(), proto.Expected(), proto.URL(), "")
		require.NoError(t, err)

		job, err := repo.Create(ctx, proto)
		require.NoError(t, err)

		if i%3 == 0 {
			err := job.Run()
			require.NoError(t, err)
			err = repo.Update(context.Background(), job)
			require.NoError(t, err)

		}
	}

	runningJobsB, err := repo.UserJobsWithState(ctx, userB, scripts.JobRunning)
	require.NoError(t, err)
	require.NotNil(t, runningJobsB)
	require.NotEmpty(t, runningJobsB)

	for _, job := range runningJobsB {
		require.Equal(t, userB, job.OwnerID())
		require.Equal(t, scripts.JobRunning, job.State())
	}

	runningJobsA, err := repo.UserJobsWithState(ctx, userA, scripts.JobRunning)
	require.NoError(t, err)
	require.NotNil(t, runningJobsA)
	require.NotEmpty(t, runningJobsB)

	for _, job := range runningJobsA {
		require.Equal(t, userA, job.OwnerID())
		require.Equal(t, scripts.JobRunning, job.State())
	}
}

func generateRandomJobPrototype(t *testing.T) *scripts.JobPrototype {
	t.Helper()

	ownerID := scripts.UserID(gofakeit.IntRange(1, 1000))
	scriptID := scripts.ScriptID(gofakeit.IntRange(1, 3))

	input := generateRandomValuePrototype(t)

	val, err := scripts.NewValueType("integer")
	require.NoError(t, err)

	outputP1, err := scripts.NewField(*val, gofakeit.LetterN(10), gofakeit.LetterN(10), gofakeit.LetterN(10))
	require.NoError(t, err)

	proto, err := scripts.NewJobPrototype(ownerID, scriptID, input, []scripts.Field{*outputP1}, "example", "")
	require.NoError(t, err)

	return proto
}

func generateRandomValuePrototype(t *testing.T) []scripts.Value {
	t.Helper()
	arr := make([]scripts.Value, 0, 2)
	for i := 0; i < cap(arr); i++ {
		numStr := strconv.Itoa(gofakeit.Number(0, 100))
		v, err := scripts.NewValue("integer", numStr)
		require.NoError(t, err)
		arr = append(arr, v)
	}
	return arr
}
