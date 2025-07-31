package service_test

import (
	"testing"

	"github.com/bmstu-itstech/scriptum-back/internal/service"
	"github.com/stretchr/testify/require"
)

func setUpMockJobtRepository() (*service.MockJobRepo, error) {
	return service.NewMockJobRepository()
}
func TestMockJobRepository(t *testing.T) {
	r, err := setUpMockJobtRepository()
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
