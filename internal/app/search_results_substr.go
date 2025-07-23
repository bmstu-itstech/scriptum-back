package app

import (
	"context"
	"log/slog"

	"github.com/bmstu-itstech/scriptum-back/internal/domain/scripts"
)

type SearchResultsSubstrUC struct {
	jobR    scripts.JobRepository
	resR    scripts.ResultRepository
	userR   scripts.UserRepository
	scriptR scripts.ScriptRepository
	logger  *slog.Logger
}

func NewSearchResultsSubstrUC(
	jobR scripts.JobRepository,
	resR scripts.ResultRepository,
	userR scripts.UserRepository,
	scriptR scripts.ScriptRepository,
	logger *slog.Logger,
) SearchResultsSubstrUC {
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
	if logger == nil {
		panic(scripts.ErrInvalidLogger)
	}
	return SearchResultsSubstrUC{
		jobR:    jobR,
		userR:   userR,
		scriptR: scriptR,
		resR:    resR,
		logger:  logger,
	}
}

func (u *SearchResultsSubstrUC) SearchResultBySubstr(ctx context.Context, userID uint32, substr string) ([]ResultDTO, error) {
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
