package usecase

import (
	"context"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type SearchJobResultsUC struct {
	jobR    scripts.JobRepository
	resR    scripts.ResultRepository
	userR   scripts.UserRepository
	scriptR scripts.ScriptRepository
}

func NewSearchJobResultsUC(
	jobR scripts.JobRepository,
	resR scripts.ResultRepository,
	userR scripts.UserRepository,
	scriptR scripts.ScriptRepository,
) (*SearchJobResultsUC, error) {
	if jobR == nil {
		panic(scripts.ErrInvalidJobRepository)
	}
	if resR == nil {
		panic(scripts.ErrInvalidResultRepository)
	}
	if userR == nil {
		panic(scripts.ErrInvalidUserRepository)
	}
	if scriptR == nil {
		panic(scripts.ErrInvalidScriptRepository)
	}
	return &SearchJobResultsUC{
		jobR:    jobR,
		userR:   userR,
		scriptR: scriptR,
		resR:    resR,
	}, nil
}

func (u *SearchJobResultsUC) SearchJobResults(ctx context.Context, userID uint32, substr string) ([]ResultDTO, error) {
	results, err := u.resR.SearchResult(ctx, userID, substr)
	if err != nil {
		return nil, err
	}

	dto := make([]ResultDTO, 0, len(results))
	for _, res := range results {
		dto = append(dto, ResultToDTO(res))
	}

	return dto, nil
}
