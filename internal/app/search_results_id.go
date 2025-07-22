package app

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type SearchResultIDUC struct {
	jobR    scripts.JobRepository
	resR    scripts.ResultRepository
	userR   scripts.UserRepository
	scriptR scripts.ScriptRepository
}

func NewSearchResultIDUCUC(
	jobR scripts.JobRepository,
	resR scripts.ResultRepository,
	userR scripts.UserRepository,
	scriptR scripts.ScriptRepository,
) (*SearchResultIDUC, error) {
	if jobR == nil {
		return nil, scripts.ErrInvalidJobRepository
	}
	if resR == nil {
		return nil, scripts.ErrInvalidResultRepository
	}
	if userR == nil {
		return nil, scripts.ErrInvalidUserRepository
	}
	if scriptR == nil {
		return nil, scripts.ErrInvalidScriptRepository
	}
	return &SearchResultIDUC{
		jobR:    jobR,
		userR:   userR,
		scriptR: scriptR,
		resR:    resR,
	}, nil
}

func (u *SearchResultIDUC) SearchResultByID(ctx context.Context, JobID scripts.JobID) (ResultDTO, error) {
	result, err := u.resR.JobResult(ctx, JobID)
	if err != nil {
		return ResultDTO{}, err
	}

	return ResultToDTO(result), nil
}
