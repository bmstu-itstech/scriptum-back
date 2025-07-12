package service

import (
	"fmt"
	"sync"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type MockJobRepo struct {
	sync.RWMutex
	m map[scripts.JobID]scripts.Job
}

func NewMockJobRepository() (*MockJobRepo, error) {
	return &MockJobRepo{
		m: make(map[scripts.JobID]scripts.Job),
	}, nil
}

func (r *MockJobRepo) GetJob(JobID scripts.JobID) (scripts.Job, error) {
	r.RLock()
	defer r.RUnlock()

	job, ok := r.m[JobID]
	if !ok {
		return scripts.Job{}, fmt.Errorf("user not found: %d", JobID)
	}
	return job, nil
}
