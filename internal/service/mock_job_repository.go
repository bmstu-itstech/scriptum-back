package service

import (
	"context"
	"fmt"
	"sync"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type MockJobRepo struct {
	sync.RWMutex
	m      map[scripts.JobID]scripts.Job
	lastID scripts.JobID
}

func NewMockJobRepository() (*MockJobRepo, error) {
	return &MockJobRepo{
		m: make(map[scripts.JobID]scripts.Job),
	}, nil
}

func (r *MockJobRepo) Create(ctx context.Context, job *scripts.JobPrototype) (*scripts.Job, error) {
	r.Lock()
	defer r.Unlock()

	r.lastID++

	newScript, err := job.Build(r.lastID)
	if err != nil {
		return nil, err
	}

	r.m[r.lastID] = *newScript

	return newScript, nil
}

func (r *MockJobRepo) Update(ctx context.Context, job *scripts.Job) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[job.ID()]; !ok {
		return fmt.Errorf("%w Update: cannot update job with id: %d", scripts.ErrJobNotFound, job.ID())
	}
	r.m[job.ID()] = *job
	return nil
}

func (r *MockJobRepo) Delete(ctx context.Context, jobID scripts.JobID) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.m[jobID]; !ok {
		return fmt.Errorf("%w Delete: cannot delete job with id: %d", scripts.ErrJobNotFound, jobID)
	}
	delete(r.m, jobID)
	return nil
}

func (r *MockJobRepo) Job(ctx context.Context, jobID scripts.JobID) (*scripts.Job, error) {
	r.RLock()
	defer r.RUnlock()

	job, ok := r.m[jobID]
	if !ok {
		return &scripts.Job{}, fmt.Errorf("%w Job: cannot extract job with id: %d", scripts.ErrJobNotFound, jobID)
	}
	return &job, nil
}

func (r *MockJobRepo) UserJobs(ctx context.Context, userID scripts.UserID) ([]scripts.Job, error) {
	r.RLock()
	defer r.RUnlock()
	jobArray := make([]scripts.Job, 0, len(r.m))
	for _, j := range r.m {
		if j.OwnerID() == userID {
			jobArray = append(jobArray, j)
		}
	}
	if len(jobArray) == 0 {
		return nil, nil
	}
	return jobArray, nil
}

func (r *MockJobRepo) UserJobsWithState(ctx context.Context, userID scripts.UserID, jobState scripts.JobState) ([]scripts.Job, error) {
	r.RLock()
	defer r.RUnlock()
	jobArray := make([]scripts.Job, 0, len(r.m))
	for _, j := range r.m {
		if j.OwnerID() == userID && j.State() == jobState {
			jobArray = append(jobArray, j)
		}
	}
	if len(jobArray) == 0 {
		return nil, nil
	}
	return jobArray, nil
}
