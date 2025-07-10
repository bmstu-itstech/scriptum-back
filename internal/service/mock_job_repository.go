package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type MockJobRepository interface {
	GetJob(JobID scripts.JobID) (scripts.Job, error)
}

type MockJobRepo struct {
	context context.Context
	sync.RWMutex
	m map[scripts.JobID]scripts.Job
}

func NewMockJobRepository() *MockJobRepo {
	return &MockJobRepo{
		context: context.Background(),
		m:       make(map[scripts.JobID]scripts.Job),
	}
}

func (r *MockJobRepo) GetUser(JobID scripts.JobID) (scripts.Job, error) {
	r.RLock()
	defer r.RUnlock()

	job, ok := r.m[JobID]
	if !ok {
		return scripts.Job{}, fmt.Errorf("user not found: %d", JobID)
	}
	return job, nil
}
